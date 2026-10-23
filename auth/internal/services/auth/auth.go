package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/auth/internal/models"
	"github.com/vald3mare/Rudstock/auth/internal/tokens"
)

// AuthService регистрация, вход и выдача токенов
type AuthService interface {
	// Register создаёт покупателя (роль всегда customer) и сразу выдаёт токены.
	// User в результате нужен для ответа {user: {id, email, role}}.
	// Ошибки: models.ErrInvalidEmail, models.ErrWeakPassword, models.ErrEmailTaken.
	Register(ctx context.Context, in RegUserInput) (AuthResult, error)

	// Login проверяет email и пароль и выдаёт токены; ответ как у Register.
	// Неверный email и неверный пароль дают одну и ту же ошибку, заблокированный отдельную.
	// Ошибки: models.ErrInvalidCredentials, models.ErrUserInactive.
	Login(ctx context.Context, in LoginUserInput) (AuthResult, error)

	// Refresh меняет refresh-токен на новую пару; старый refresh после этого недействителен.
	// Повторное использование уже обменянного refresh гасит все сессии этого входа.
	// Ошибки: models.ErrInvalidRefresh.
	Refresh(ctx context.Context, refresh string) (Tokens, error)

	// Logout гасит вход, которому принадлежит refresh. Неизвестный или погашенный refresh не ошибка.
	Logout(ctx context.Context, refresh string) error
}

// это интерфейс сущности которая заведует токенами (пакет tokens)
type Emitter interface {
	Emit(userID uuid.UUID, role models.Role, sessionID uuid.UUID) (string, error)
	Verify(access string) (tokens.Identity, error)
	Refresh() (raw string, hash []byte, err error)
	HashRefresh(raw string) []byte
}

// SessionRepo что сервису нужно от хранилища сессий
type SessionRepo interface {
	// Create сохраняет новую сессию и возвращает её id.
	Create(ctx context.Context, s models.Session) (uuid.UUID, error)

	// GetByTokenHash возвращает сессию по хешу refresh, в том числе погашенную.
	// Ошибки: models.ErrSessionNotFound.
	GetByTokenHash(ctx context.Context, hash []byte) (models.Session, error)

	// Rotate гасит старую сессию и создаёт следующую одной транзакцией.
	// Ошибки: models.ErrSessionRevoked, если старую уже погасили.
	Rotate(ctx context.Context, oldID uuid.UUID, next models.Session) (uuid.UUID, error)

	// RevokeFamily гасит все живые сессии одного входа.
	RevokeFamily(ctx context.Context, familyID uuid.UUID) error
}

// UserRepo что сервису нужно от хранилища пользователей.
type UserRepo interface {
	// Create сохраняет пользователя и возвращает сгенерированный id.
	// Ошибки: models.ErrEmailTaken.
	Create(ctx context.Context, user models.User) (uuid.UUID, error)

	// GetByEmail возвращает пользователя по email (для Login).
	// Ошибки: models.ErrUserNotFound.
	GetByEmail(ctx context.Context, email string) (models.User, error)

	// GetByID возвращает пользователя по id (для Refresh).
	// Ошибки: models.ErrUserNotFound.
	GetByID(ctx context.Context, id uuid.UUID) (models.User, error)
}
type authService struct {
	userRepo    UserRepo
	sessionRepo SessionRepo
	emitter     Emitter

	refreshTTL    time.Duration // скользящий срок: каждый refresh продлевает вход на столько
	refreshMaxTTL time.Duration // абсолютный предел от ввода пароля, дальше только новый login
}

func NewAuthService(users UserRepo, sessions SessionRepo, emitter Emitter, refreshTTL, refreshMaxTTL time.Duration) AuthService {
	return &authService{
		userRepo:      users,
		sessionRepo:   sessions,
		emitter:       emitter,
		refreshTTL:    refreshTTL,
		refreshMaxTTL: refreshMaxTTL,
	}
}
