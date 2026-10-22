package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const revokeFamilyQuery = `
		UPDATE sessions
		SET revoked_at = now()
		WHERE family_id = $1 AND revoked_at IS NULL`

// RevokeFamily гасит все живые сессии семьи. 0 затронутых строк не ошибка: семья могла быть уже погашена.
func (r *sessionRepo) RevokeFamily(ctx context.Context, familyID uuid.UUID) error {
	if _, err := r.pool.Exec(ctx, revokeFamilyQuery, familyID); err != nil {
		return fmt.Errorf("revoke session family: %w", err)
	}

	return nil
}
