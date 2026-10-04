package category

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// CategoryRepo реализация service.CategoryRepo поверх Postgres.
type CategoryRepo struct {
	pool *pgxpool.Pool
}

// NewCategoryRepo конструктор репозитория категорий.
func NewCategoryRepo(pool *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{
		pool: pool,
	}
}
