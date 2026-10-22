package auth

import (
	"unicode"

	"github.com/alexedwards/argon2id"
)

// dummyHash для Login с несуществующим email: сравниваем пароль с ним, чтобы ответ шёл столько же,
// сколько с существующим, и по времени нельзя было узнать, зарегистрирован ли email
var dummyHash, _ = argon2id.CreateHash("dummy-password-for-timing", argon2id.DefaultParams)

func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasLetter, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	return hasLetter && hasDigit && hasSpecial
}
