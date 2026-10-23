package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const revokeAllByUserQuery = `
		UPDATE sessions
		SET revoked_at = now()
		WHERE user_id = $1 AND revoked_at IS NULL`

// RevokeAllByUser гасит все живые сессии пользователя. 0 затронутых строк не ошибка: живых могло не быть.
func (r *sessionRepo) RevokeAllByUser(ctx context.Context, userID uuid.UUID) error {
	if _, err := r.pool.Exec(ctx, revokeAllByUserQuery, userID); err != nil {
		return fmt.Errorf("revoke user sessions: %w", err)
	}

	return nil
}
