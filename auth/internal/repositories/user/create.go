package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

const createQuery = `
		INSERT INTO users (email, password_hash, role, is_active, email_verified)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

func (r *userRepo) Create(ctx context.Context, user models.User) (uuid.UUID, error) {
	var userID uuid.UUID
	err := r.pool.QueryRow(ctx, createQuery,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.IsActive,
		user.EmailVerified,
	).Scan(&userID)
	if err != nil {
		// 23505 unique_violation по email; имя ограничения проверяем, чтобы будущие UNIQUE не выдавались за "email занят"
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
			return uuid.Nil, models.ErrEmailTaken
		}
		return uuid.Nil, fmt.Errorf("insert user: %w", err)
	}

	return userID, nil
}
