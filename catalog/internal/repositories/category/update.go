package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// тот же приём, что у карточки: nil уходит как NULL, COALESCE оставляет старое значение
const updateQuery = `
		UPDATE categories
		SET name = COALESCE($2, name)
		WHERE id = $1
		RETURNING id, name`

// Update применяет patch одним запросом и возвращает категорию после изменения.
// Если категории нет, возвращает models.ErrCategoryNotFound, если имя занято, models.ErrCategoryExists.
func (r *categoryRepo) Update(ctx context.Context, id int64, patch models.CategoryPatch) (models.Category, error) {
	var c models.Category
	err := r.pool.QueryRow(ctx, updateQuery, id, patch.Name).Scan(&c.ID, &c.Name)
	if err != nil {
		// UPDATE без подходящей строки ничего не возвращает: категории с таким id нет
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Category{}, models.ErrCategoryNotFound
		}
		// 23505 unique_violation: другая категория уже с таким name
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.Category{}, models.ErrCategoryExists
		}
		return models.Category{}, fmt.Errorf("update category: %w", err)
	}

	return c, nil
}
