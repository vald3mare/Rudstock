package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

// revoked_at и expires_at не фильтруем: решение по погашенной/истёкшей сессии принимает сервис
const getByTokenHashQuery = `
		SELECT ` + sessionColumns + `
		FROM sessions
		WHERE token_hash = $1`

// GetByTokenHash читает сессию по хешу refresh-токена. Если строки нет, возвращает models.ErrSessionNotFound.
func (r *sessionRepo) GetByTokenHash(ctx context.Context, hash []byte) (models.Session, error) {
	s, err := scanSession(r.pool.QueryRow(ctx, getByTokenHashQuery, hash))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Session{}, models.ErrSessionNotFound
		}
		return models.Session{}, fmt.Errorf("select session by token hash: %w", err)
	}

	return s, nil
}
