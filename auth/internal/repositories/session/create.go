package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

// id и created_at не передаем, их проставляют DEFAULT'ы; revoked_at NULL = живая
const createQuery = `
		INSERT INTO sessions (user_id, family_id, token_hash, authenticated_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

func (r *sessionRepo) Create(ctx context.Context, s models.Session) (uuid.UUID, error) {
	return insert(ctx, r.pool, s)
}

func insert(ctx context.Context, q querier, s models.Session) (uuid.UUID, error) {
	var sessionID uuid.UUID
	err := q.QueryRow(ctx, createQuery,
		s.UserID,
		s.FamilyID,
		s.TokenHash,
		s.AuthenticatedAt,
		s.ExpiresAt,
	).Scan(&sessionID)
	if err != nil {
		// 23503 foreign_key_violation: пользователя с таким id нет
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return uuid.Nil, models.ErrUserNotFound
		}
		return uuid.Nil, fmt.Errorf("insert session: %w", err)
	}

	return sessionID, nil
}
