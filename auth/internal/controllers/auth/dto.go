package auth

import (
	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/auth/internal/models"
	authsvc "github.com/vald3mare/Rudstock/auth/internal/services/auth"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshRequest тело /auth/refresh и /auth/logout
type RefreshRequest struct {
	Refresh string `json:"refresh"`
}

type TokensResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

// UserResponse только то, что в контракте: PasswordHash и прочее наружу не уходят
type UserResponse struct {
	ID    uuid.UUID   `json:"id"`
	Email string      `json:"email"`
	Role  models.Role `json:"role"`
}

// AuthResponse ответ register/login: {access, refresh, user: {id, email, role}}
type AuthResponse struct {
	Access  string       `json:"access"`
	Refresh string       `json:"refresh"`
	User    UserResponse `json:"user"`
}

func toTokensResponse(t authsvc.Tokens) TokensResponse {
	return TokensResponse{
		Access:  t.Access,
		Refresh: t.Refresh,
	}
}

func toAuthResponse(res authsvc.AuthResult) AuthResponse {
	return AuthResponse{
		Access:  res.Tokens.Access,
		Refresh: res.Tokens.Refresh,
		User: UserResponse{
			ID:    res.User.ID,
			Email: res.User.Email,
			Role:  res.User.Role,
		},
	}
}

func (r RegisterRequest) toInput() authsvc.RegUserInput {
	return authsvc.RegUserInput{
		Email:    r.Email,
		Password: r.Password,
	}
}

func (r LoginRequest) toInput() authsvc.LoginUserInput {
	return authsvc.LoginUserInput{
		Email:    r.Email,
		Password: r.Password,
	}
}
