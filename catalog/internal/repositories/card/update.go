package card

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// nil-указатель уходит в базу как NULL, а COALESCE(NULL, колонка) оставляет старое значение.
// Так не нужно собирать SET строкой под каждый набор полей.
// RETURNING отдаёт карточку целиком, второй SELECT не нужен.
const updateQuery = `
		UPDATE cards
		SET category_id = COALESCE($2, category_id),
		    title       = COALESCE($3, title),
		    description = COALESCE($4, description),
		    price       = COALESCE($5, price),
		    photo_url   = COALESCE($6, photo_url)
		WHERE id = $1
		RETURNING ` + cardColumns

// Update применяет patch одним запросом и возвращает карточку после изменения.
// Если карточки нет, возвращает models.ErrCardNotFound, если новой категории нет, models.ErrCategoryNotFound.
func (r *cardRepo) Update(ctx context.Context, id uuid.UUID, patch models.CardPatch) (models.Card, error) {
	c, err := scanCard(r.pool.QueryRow(ctx, updateQuery,
		id,
		patch.CategoryID,
		patch.Title,
		patch.Description,
		patch.Price,
		patch.PhotoURL,
	))
	if err != nil {
		// UPDATE без подходящей строки ничего не возвращает: карточки с таким id нет
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Card{}, models.ErrCardNotFound
		}
		// 23503 foreign_key_violation: новой категории с таким id нет
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return models.Card{}, models.ErrCategoryNotFound
		}
		return models.Card{}, fmt.Errorf("update card: %w", err)
	}

	return c, nil
}
