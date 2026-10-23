package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

// Rotate гасит старую сессию и создаёт следующую одной транзакцией:
// без неё при падении между запросами клиент остался бы либо без сессии, либо с двумя живыми.
func (r *sessionRepo) Rotate(ctx context.Context, oldID uuid.UUID, next models.Session) (uuid.UUID, error) {
	var nextID uuid.UUID
	err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, revokeQuery, oldID)
		if err != nil {
			return fmt.Errorf("revoke session: %w", err)
		}
		// старую уже погасил параллельный refresh: новую не создаём, транзакция откатится
		if tag.RowsAffected() == 0 {
			return models.ErrSessionRevoked
		}

		nextID, err = insert(ctx, tx, next)
		return err
	})
	if err != nil {
		return uuid.Nil, err
	}

	return nextID, nil
}
