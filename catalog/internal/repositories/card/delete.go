package card

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// Delete удаляет карточку по id. Если строки нет, возвращает models.ErrCardNotFound.
func (r *CardRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM cards
		WHERE id = $1`

	// Exec, а не QueryRow: DELETE без RETURNING строк не возвращает,
	// поэтому "не найдено" узнаём по числу затронутых строк, а не по ErrNoRows
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete card: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrCardNotFound
	}

	return nil
}
