package card

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// id и created_at не передаем, их проставляют DEFAULT'ы в таблице
const creteQuery = `
		INSERT INTO cards (category_id, title, description, price, photo_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

// вот и сама реализация создания карточки, именно от сюда мы и ходим в бд
func (r *CardRepo) Create(ctx context.Context, card models.Card) (uuid.UUID, error) {
	var card_id uuid.UUID
	err := r.pool.QueryRow(ctx, creteQuery,
		card.CategoryID,
		card.Title,
		card.Description,
		card.Price,
		card.PhotoURL,
	).Scan(&card_id)
	if err != nil {
		// 23503 foreign_key_violation: категории с таким id нет
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return uuid.Nil, models.ErrCategoryNotFound
		}
		return uuid.Nil, fmt.Errorf("insert card: %w", err)
	}

	return card_id, nil
}
