package auth

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

func (s *authService) Register(ctx context.Context, in RegUserInput) (AuthResult, error) {
	fmtEmail := strings.ToLower(strings.TrimSpace(in.Email))
	addr, err := mail.ParseAddress(fmtEmail)
	if err != nil || addr.Address != fmtEmail {
		return AuthResult{}, models.ErrInvalidEmail
	}

	if len(in.Password) < 8 {
		return AuthResult{}, models.ErrWeakPassword
	}

	hash, err := argon2id.CreateHash(in.Password, argon2id.DefaultParams)
	if err != nil {
		return AuthResult{}, fmt.Errorf("hash password: %w", err)
	}

	user := models.User{
		Email:        fmtEmail,
		PasswordHash: hash,
		Role:         models.RoleCustomer,
		IsActive:     true,
	}

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return AuthResult{}, err
	}

	//TODO выпуск токена
	user.ID = id
	return AuthResult{User: user}, nil
}
