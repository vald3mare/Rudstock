package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

// newSession генерирует refresh-токен и собирает под него сессию; в базу не пишет.
// authenticatedAt при ротации переносится из старой сессии, от него считается абсолютный предел.
func (s *authService) newSession(userID, familyID uuid.UUID, authenticatedAt time.Time) (string, models.Session, error) {
	raw, hash, err := s.emitter.Refresh()
	if err != nil {
		return "", models.Session{}, fmt.Errorf("generate refresh: %w", err)
	}

	// ExpiresAt = min(скользящий от сейчас, абсолютный от ввода пароля)
	expiresAt := time.Now().Add(s.refreshTTL)
	if limit := authenticatedAt.Add(s.refreshMaxTTL); limit.Before(expiresAt) {
		expiresAt = limit
	}

	return raw, models.Session{
		UserID:          userID,
		FamilyID:        familyID,
		TokenHash:       hash,
		AuthenticatedAt: authenticatedAt,
		ExpiresAt:       expiresAt,
	}, nil
}

// startSession открывает новый вход (register/login): новая семья, сессия в базе, пара токенов.
func (s *authService) startSession(ctx context.Context, user models.User) (Tokens, error) {
	raw, session, err := s.newSession(user.ID, uuid.New(), time.Now())
	if err != nil {
		return Tokens{}, err
	}

	// сначала сессия: её id нужен в access-токене как sid
	sessionID, err := s.sessionRepo.Create(ctx, session)
	if err != nil {
		return Tokens{}, err
	}

	access, err := s.emitter.Emit(user.ID, user.Role, sessionID)
	if err != nil {
		return Tokens{}, fmt.Errorf("emit access: %w", err)
	}

	return Tokens{Access: access, Refresh: raw}, nil
}
