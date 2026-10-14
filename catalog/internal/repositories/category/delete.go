package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

const deleteQuery = `
		DELETE FROM categories
		WHERE id = $1`

// Delete удаляет категорию по id.
// Если строки нет, возвращает models.ErrCategoryNotFound, если в ней есть карточки, models.ErrCategoryHasCards.
func (r *categoryRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, deleteQuery, id)
	if err != nil {
		// 23503 foreign_key_violation: на категорию ссылаются карточки (cards.category_id)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return models.ErrCategoryHasCards
		}
		return fmt.Errorf("delete category: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrCategoryNotFound
	}

	return nil
}
