package category

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

const listQuery = `
		SELECT id, name
		FROM categories
		ORDER BY name`

// List читает все категории, отсортированные по имени.
func (r *categoryRepo) List(ctx context.Context) ([]models.Category, error) {
	rows, err := r.pool.Query(ctx, listQuery)
	if err != nil {
		return nil, fmt.Errorf("select categories: %w", err)
	}

	// CollectRows сам закрывает rows и возвращает пустой срез, а не nil, если строк нет
	categories, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Category, error) {
		var c models.Category
		err := row.Scan(&c.ID, &c.Name)
		return c, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan categories: %w", err)
	}

	return categories, nil
}
