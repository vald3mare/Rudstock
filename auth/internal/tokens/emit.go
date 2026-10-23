package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

// Emit выпускает подписанный access-токен.
func (t *token) Emit(userID uuid.UUID, role models.Role, sessionID uuid.UUID) (string, error) {
	now := time.Now()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		Role:      role,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    t.publisher,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(t.ttl)),
		},
	})

	return jwtToken.SignedString(t.secret)
}
