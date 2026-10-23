package tokens

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Verify проверяет подпись, срок и издателя access-токена.
// Ошибки: ErrTokenExpired (клиенту идти в /auth/refresh), ErrInvalidToken (всё остальное).
func (t *token) Verify(access string) (Identity, error) {
	var c claims
	_, err := jwt.ParseWithClaims(access, &c, func(*jwt.Token) (any, error) {
		return t.secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), //запрещаем другие алгоритмы шифрования
		jwt.WithIssuer(t.publisher),
		jwt.WithExpirationRequired(),
	)
	if errors.Is(err, jwt.ErrTokenExpired) {
		return Identity{}, ErrTokenExpired
	}
	if err != nil {
		return Identity{}, ErrInvalidToken
	}

	userID, err := uuid.Parse(c.Subject)
	if err != nil {
		return Identity{}, ErrInvalidToken
	}

	return Identity{
		UserID:    userID,
		Role:      c.Role,
		SessionID: c.SessionID,
	}, nil
}
