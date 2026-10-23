package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

func (s *authService) Login(ctx context.Context, in LoginUserInput) (AuthResult, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if errors.Is(err, models.ErrUserNotFound) {
		// результат не важен, важно потраченное время (см. dummyHash)
		_, _ = argon2id.ComparePasswordAndHash(in.Password, dummyHash)
		return AuthResult{}, models.ErrInvalidCredentials
	}
	if err != nil {
		return AuthResult{}, err
	}

	ok, err := argon2id.ComparePasswordAndHash(in.Password, user.PasswordHash)
	if err != nil {
		return AuthResult{}, fmt.Errorf("compare password: %w", err)
	}
	if !ok {
		return AuthResult{}, models.ErrInvalidCredentials
	}

	// блокировку проверяем только после верного пароля: иначе 403 выдаёт, что такой email есть
	if !user.IsActive {
		return AuthResult{}, models.ErrUserInactive
	}

	tokens, err := s.startSession(ctx, user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{Tokens: tokens, User: user}, nil
}
