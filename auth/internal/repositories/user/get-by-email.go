package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vald3mare/Rudstock/auth/internal/models"
)

// email ищем как есть: в нижний регистр его приводит сервис, CHECK в таблице гарантирует, что в базе он такой же
const getByEmailQuery = `
		SELECT ` + userColumns + `
		FROM users
		WHERE email = $1`

// GetByEmail читает пользователя по email. Если строки нет, возвращает models.ErrUserNotFound.
func (r *userRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	u, err := scanUser(r.pool.QueryRow(ctx, getByEmailQuery, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("select user by email: %w", err)
	}

	return u, nil
}
