package auth

import (
	"context"
	"errors"

	"github.com/vald3mare/Rudstock/auth/internal/models"
)

// Logout гасит всю семью, а не одну сессию: если клиент прислал уже обменянный refresh
func (s *authService) Logout(ctx context.Context, refresh string) error {
	session, err := s.sessionRepo.GetByTokenHash(ctx, s.emitter.HashRefresh(refresh))
	if errors.Is(err, models.ErrSessionNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	return s.sessionRepo.RevokeFamily(ctx, session.FamilyID)
}
