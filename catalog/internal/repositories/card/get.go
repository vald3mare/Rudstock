package card

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

const getQuery = `
		SELECT ` + cardColumns + `
		FROM cards
		WHERE id = $1`

// Get читает карточку по id. Если строки нет, возвращает models.ErrCardNotFound.
func (r *cardRepo) Get(ctx context.Context, id uuid.UUID) (models.Card, error) {
	c, err := scanCard(r.pool.QueryRow(ctx, getQuery, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Card{}, models.ErrCardNotFound
		}
		return models.Card{}, fmt.Errorf("select card: %w", err)
	}

	return c, nil
}
