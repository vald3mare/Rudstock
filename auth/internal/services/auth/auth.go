package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

//   "auth": {
//     "POST /auth/register": { "status": "todo", "request": ["email", "password"], "response": ["access", "refresh", "user: {id, email, role}"], "errors": [400, 409, 429], "note": "роль всегда customer; name и phone не спрашиваем, заполняются в PATCH /me или при первом заказе" },
//     "POST /auth/login": { "status": "todo", "request": ["email", "password"], "response": ["access", "refresh", "user: {id, email, role}"], "errors": [401, 403, 429] },
//     "POST /auth/refresh": { "status": "todo", "request": ["refresh"], "response": ["access", "refresh"], "errors": [401] },
//     "POST /auth/logout": { "status": "todo", "request": ["refresh"] },
//     "POST /auth/email/verify": { "status": "todo", "request": ["token"], "errors": [400] },
//     "POST /auth/password/forgot": { "status": "todo", "request": ["email"], "errors": [429], "note": "204 всегда, даже для несуществующего email" },
//     "POST /auth/password/reset": { "status": "todo", "request": ["token", "password"], "errors": [400] }
//   },

// AuthService регистрация, вход и выдача токенов.
type AuthService interface {
	// Register создаёт покупателя (роль всегда customer) и сразу выдаёт токены.
	// User в результате нужен для ответа {user: {id, email, role}}.
	// Ошибки: models.ErrInvalidEmail, models.ErrWeakPassword, models.ErrEmailTaken.
	Register(ctx context.Context, in RegUserInput) (AuthResult, error)

	// Login проверяет email и пароль и выдаёт токены; ответ как у Register.
	// Неверный email и неверный пароль дают одну и ту же ошибку, заблокированный отдельную.
	Login(ctx context.Context, in LoginUserInput) (AuthResult, error)
}

// UserRepo что сервису нужно от хранилища пользователей.
type UserRepo interface {
	// Create сохраняет пользователя и возвращает сгенерированный id.
	// Ошибки: email занят.
	Create(ctx context.Context, user models.User) (uuid.UUID, error)

	// GetByEmail возвращает пользователя по email (для Login).
	// Ошибки: пользователь не найден.
	GetByEmail(ctx context.Context, email string) (models.User, error)
}

type authService struct {
	userRepo UserRepo
}

func NewAuthService(repo UserRepo) AuthService {
	return &authService{userRepo: repo}
}
