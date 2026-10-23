package models

import (
	"time"

	"github.com/google/uuid"
)

// Session одна выдача refresh-токена. При ротации старая сессия гасится, создаётся новая с тем же FamilyID.
// Сам токен не хранится, только его хеш: утечка БД не должна давать готовые токены.
type Session struct {
	ID       uuid.UUID // попадает в access-токен как sid
	UserID   uuid.UUID
	FamilyID uuid.UUID // общий для цепочки login -> refresh -> refresh; повторное использование гасит всю семью

	TokenHash []byte // sha256 от refresh-токена, по нему ищем сессию

	// AuthenticatedAt момент ввода пароля; при ротации переносится как есть,
	// от него считается абсолютный предел жизни семьи.
	AuthenticatedAt time.Time

	// ExpiresAt итоговый срок: min(скользящий от CreatedAt, абсолютный от AuthenticatedAt).
	// Считается сервисом при создании, при refresh проверяется только он.
	ExpiresAt time.Time
	CreatedAt time.Time

	RevokedAt *time.Time // nil = живая; иначе использована при ротации или отозвана (logout, сброс пароля)
}
