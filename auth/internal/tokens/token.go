package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

type token struct {
	secret    []byte
	publisher string
	ttl       time.Duration
}

// то что мы достаем из Verify
type Identity struct {
	UserID    uuid.UUID
	Role      models.Role
	SessionID uuid.UUID
}

type claims struct {
	Role      models.Role `json:"role"`
	SessionID uuid.UUID   `json:"sid"`
	jwt.RegisteredClaims
}

// сущность по управлению токенами
type Emitter interface {
	Emit(userID uuid.UUID, role models.Role, sessionID uuid.UUID) (string, error)
	Verify(access string) (Identity, error)
	Refresh() (raw string, hash []byte, err error)
	HashRefresh(raw string) []byte
}

const minSecretSize = 32

func NewEmitter(secret []byte, publisher string, ttl time.Duration) (Emitter, error) {
	if len(secret) < minSecretSize {
		return nil, fmt.Errorf("jwt secret must be at least %d bytes", minSecretSize)
	}
	if publisher == "" {
		return nil, errors.New("jwt publisher is empty")
	}
	if ttl <= 0 {
		return nil, errors.New("access token ttl must be positive")
	}

	return &token{
		secret:    secret,
		publisher: publisher,
		ttl:       ttl,
	}, nil
}
