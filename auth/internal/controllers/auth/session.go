package auth

import (
	"context"

	authsvc "github.com/vald3mare/Rudstock/auth/internal/services/auth"
)

// sessionService работа с уже открытой сессией по refresh-токену.
type sessionService interface {
	Refresh(ctx context.Context, refresh string) (authsvc.Tokens, error)
	Logout(ctx context.Context, refresh string) error
}

type SessionHandler struct {
	svc sessionService
}

func NewSessionHandler(s sessionService) *SessionHandler {
	return &SessionHandler{svc: s}
}
