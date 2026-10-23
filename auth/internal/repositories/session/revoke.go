package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

// revoked_at IS NULL: гасим только живую; второй такой же запрос затронет 0 строк,
// так два параллельных refresh одним токеном не пройдут оба
const revokeQuery = `
		UPDATE sessions
		SET revoked_at = now()
		WHERE id = $1 AND revoked_at IS NULL`

// Revoke гасит сессию. Если живой сессии с таким id нет, возвращает models.ErrSessionRevoked.
func (r *sessionRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, revokeQuery, id)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrSessionRevoked
	}

	return nil
}
