package services

import (
	"context"
)

// AuthService вход и токены. Профиль (имя, телефон) живёт в сервисе users,
// auth передаёт его туда при регистрации.
type AuthService interface {
	// Register создаёт покупателя (роль всегда customer) и сразу выдаёт пару токенов.
	// phone необязательный, может быть пустым.
	Register(ctx context.Context, email, password, name, phone string) (accessToken, refreshToken string, err error)

	// Login проверяет email и пароль и выдаёт пару токенов.
	// Неверный email и неверный пароль дают одну и ту же ошибку.
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)

	// Refresh гасит присланный refresh-токен и выдаёт новую пару (ротация).
	Refresh(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error)

	// Logout отзывает одну сессию по её refresh-токену.
	Logout(ctx context.Context, refreshToken string) error

	// ChangePassword меняет пароль вошедшего пользователя (userID из access-токена)
	// и отзывает все его сессии.
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error

	// ForgotPassword отправляет письмо со ссылкой сброса. Для несуществующего email
	// ошибки нет: снаружи нельзя узнать, зарегистрирован ли адрес.
	ForgotPassword(ctx context.Context, email string) error

	// ResetPassword ставит новый пароль по одноразовому токену из письма и отзывает все сессии.
	ResetPassword(ctx context.Context, resetToken, newPassword string) error

	// VerifyEmail подтверждает email по одноразовому токену из письма.
	VerifyEmail(ctx context.Context, verifyToken string) error
}
