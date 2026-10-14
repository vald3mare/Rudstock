package auth

import "github.com/vald3mare/Rudstock/auth/internal/models"

type Tokens struct {
	Access  string
	Refresh string
}

// AuthResult ответ register/login: токены и данные пользователя для {user: {...}}
type AuthResult struct {
	Tokens Tokens
	User   models.User
}

type RegUserInput struct {
	Email    string
	Password string
}

type LoginUserInput struct {
	Email    string
	Password string
}
