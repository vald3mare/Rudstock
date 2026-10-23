package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vald3mare/Rudstock/auth/internal/models"
)

func (s *authService) Refresh(ctx context.Context, refresh string) (Tokens, error) {
	session, err := s.sessionRepo.GetByTokenHash(ctx, s.emitter.HashRefresh(refresh))
	if errors.Is(err, models.ErrSessionNotFound) {
		return Tokens{}, models.ErrInvalidRefresh
	}
	if err != nil {
		return Tokens{}, err
	}

	// токен уже обменивали: его украли или это старая копия у клиента. Не знаем, у кого настоящий,
	// поэтому гасим весь вход, и обоим придётся логиниться заново
	if session.RevokedAt != nil {
		if err := s.sessionRepo.RevokeFamily(ctx, session.FamilyID); err != nil {
			return Tokens{}, err
		}
		return Tokens{}, models.ErrInvalidRefresh
	}

	if !time.Now().Before(session.ExpiresAt) {
		return Tokens{}, models.ErrInvalidRefresh
	}

	// роль и блокировку берём из базы, а не из старого токена: они могли измениться
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return Tokens{}, err
	}
	if !user.IsActive {
		if err := s.sessionRepo.RevokeFamily(ctx, session.FamilyID); err != nil {
			return Tokens{}, err
		}
		return Tokens{}, models.ErrInvalidRefresh
	}

	raw, next, err := s.newSession(user.ID, session.FamilyID, session.AuthenticatedAt)
	if err != nil {
		return Tokens{}, err
	}

	nextID, err := s.sessionRepo.Rotate(ctx, session.ID, next)
	// параллельный refresh тем же токеном успел первым: он и продолжит вход, этот отклоняем
	if errors.Is(err, models.ErrSessionRevoked) {
		return Tokens{}, models.ErrInvalidRefresh
	}
	if err != nil {
		return Tokens{}, err
	}

	access, err := s.emitter.Emit(user.ID, user.Role, nextID)
	if err != nil {
		return Tokens{}, fmt.Errorf("emit access: %w", err)
	}

	return Tokens{Access: access, Refresh: raw}, nil
}
