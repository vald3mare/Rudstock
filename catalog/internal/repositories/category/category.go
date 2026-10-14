package category

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

type CategoryRepo interface {
	Create(ctx context.Context, category models.Category) (int64, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]models.Category, error)
	Update(ctx context.Context, id int64, patch models.CategoryPatch) (models.Category, error)
}

// CategoryRepo реализация service.CategoryRepo поверх Postgres.
type categoryRepo struct {
	pool *pgxpool.Pool
}

// NewCategoryRepo конструктор репозитория категорий.
func NewCategoryRepo(pool *pgxpool.Pool) CategoryRepo {
	return &categoryRepo{
		pool: pool,
	}
}
