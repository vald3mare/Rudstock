package auth

import (
	"errors"

	"github.com/vald3mare/Rudstock/auth/internal/models"
	"github.com/vald3mare/Rudstock/auth/internal/platform/errs"
)

// mapAuthError переводит доменные ошибки в errs, чтобы httpx выбрал правильный код.
// Всё неизвестное уходит как есть, WriteError ответит 500.
func mapAuthError(err error) error {
	switch {
	case errors.Is(err, models.ErrInvalidEmail):
		return errs.InvalidInput("invalid-email", "Invalid email", err)
	case errors.Is(err, models.ErrWeakPassword):
		return errs.InvalidInput("weak-password", "Password must be at least 8 characters and contain a letter, a digit and a special character", err)
	case errors.Is(err, models.ErrEmailTaken):
		return errs.Conflict("email-taken", "Email is already registered", err)
	case errors.Is(err, models.ErrInvalidCredentials):
		return errs.Unauthorized("invalid-credentials", "Invalid email or password", err)
	case errors.Is(err, models.ErrUserInactive):
		return errs.Forbidden("user-inactive", "User is blocked", err)
	case errors.Is(err, models.ErrInvalidRefresh):
		return errs.Unauthorized("invalid-refresh", "Refresh token is invalid or expired", err)
	default:
		return err
	}
}
