package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

type UserRepo interface {
	// Create сохраняет пользователя и возвращает сгенерированный id.
	// Ошибки: email занят.
	Create(ctx context.Context, user models.User) (uuid.UUID, error)

	// GetByEmail возвращает пользователя по email (для Login).
	// Ошибки: models.ErrUserNotFound.
	GetByEmail(ctx context.Context, email string) (models.User, error)

	// GetByID возвращает пользователя по id (для Refresh: актуальная роль и блокировка).
	// Ошибки: models.ErrUserNotFound.
	GetByID(ctx context.Context, id uuid.UUID) (models.User, error)
}

type userRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo конструктор репозитория пользователей.
func NewUserRepo(pool *pgxpool.Pool) UserRepo {
	return &userRepo{pool: pool}
}

// userColumns и scanUser описывают пользователя один раз для всех SELECT/RETURNING.
const userColumns = `id, email, email_verified, password_hash, role, is_active, created_at`

func scanUser(row pgx.Row) (models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.EmailVerified, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt)
	return u, err
}
