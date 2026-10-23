package session

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

type SessionRepo interface {
	// Create сохраняет новую сессию (login/register) и возвращает сгенерированный id.
	// Ошибки: models.ErrUserNotFound.
	Create(ctx context.Context, s models.Session) (uuid.UUID, error)

	// GetByTokenHash возвращает сессию по sha256 от refresh-токена, в том числе погашенную:
	// по погашенной сервис узнаёт повторное использование токена.
	// Ошибки: models.ErrSessionNotFound.
	GetByTokenHash(ctx context.Context, hash []byte) (models.Session, error)

	// Rotate в одной транзакции гасит старую сессию и создаёт следующую в той же семье.
	// Ошибки: models.ErrSessionRevoked, если старая уже погашена (параллельный refresh тем же токеном).
	Rotate(ctx context.Context, oldID uuid.UUID, next models.Session) (uuid.UUID, error)

	// Revoke гасит одну сессию (logout).
	// Ошибки: models.ErrSessionRevoked, если живой сессии с таким id нет.
	Revoke(ctx context.Context, id uuid.UUID) error

	// RevokeFamily гасит все живые сессии семьи (повторное использование refresh-токена).
	RevokeFamily(ctx context.Context, familyID uuid.UUID) error

	// RevokeAllByUser гасит все живые сессии пользователя (сброс пароля, блокировка).
	RevokeAllByUser(ctx context.Context, userID uuid.UUID) error
}

type sessionRepo struct {
	pool *pgxpool.Pool
}

// NewSessionRepo конструктор репозитория сессий.
func NewSessionRepo(pool *pgxpool.Pool) SessionRepo {
	return &sessionRepo{pool: pool}
}

// sessionColumns и scanSession описывают сессию один раз для всех SELECT/RETURNING.
const sessionColumns = `id, user_id, family_id, token_hash, authenticated_at, expires_at, created_at, revoked_at`

func scanSession(row pgx.Row) (models.Session, error) {
	var s models.Session
	err := row.Scan(&s.ID, &s.UserID, &s.FamilyID, &s.TokenHash, &s.AuthenticatedAt, &s.ExpiresAt, &s.CreatedAt, &s.RevokedAt)
	return s, err
}

// querier общее у pool и tx: insert работает и сам по себе, и внутри Rotate.
type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
