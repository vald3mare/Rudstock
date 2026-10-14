package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// id и created_at не передаем, их проставляют DEFAULT'ы в таблице
const createQuery = `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id`

// Create вставляет категорию в categories и возвращает id, выданный базой.
func (r *categoryRepo) Create(ctx context.Context, category models.Category) (int64, error) {
	var category_id int64
	err := r.pool.QueryRow(ctx, createQuery, category.Name).Scan(&category_id)
	if err != nil {
		// 23505 unique_violation: категория с таким name уже есть
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, models.ErrCategoryExists
		}
		return 0, fmt.Errorf("insert category: %w", err)
	}

	return category_id, nil
}
