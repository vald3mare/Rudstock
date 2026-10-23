package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

const getByIDQuery = `
		SELECT ` + userColumns + `
		FROM users
		WHERE id = $1`

// GetByID читает пользователя по id. Если строки нет, возвращает models.ErrUserNotFound.
func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	u, err := scanUser(r.pool.QueryRow(ctx, getByIDQuery, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("select user by id: %w", err)
	}

	return u, nil
}
