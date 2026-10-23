package auth

import (
	"context"

	authsvc "github.com/vald3mare/Rudstock/auth/internal/services/auth"
)

// credentialsService вход по email и паролю: открывает новую сессию.
type credentialsService interface {
	Register(ctx context.Context, in authsvc.RegUserInput) (authsvc.AuthResult, error)
	Login(ctx context.Context, in authsvc.LoginUserInput) (authsvc.AuthResult, error)
}

type CredentialsHandler struct {
	svc credentialsService
}

func NewCredentialsHandler(s credentialsService) *CredentialsHandler {
	return &CredentialsHandler{svc: s}
}
