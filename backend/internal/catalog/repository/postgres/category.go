package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vald3mare/Rudstock/backend/internal/catalog/domain"
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

// Create вставляет категорию в categories и возвращает id, выданный базой.
func (r *CategoryRepo) Create(ctx context.Context, category domain.Category) (int64, error) {
	// id и created_at не передаем, их проставляют DEFAULT'ы в таблице
	const query = `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id`

	var category_id int64
	err := r.pool.QueryRow(ctx, query, category.Name).Scan(&category_id)
	if err != nil {
		// 23505 unique_violation: категория с таким name уже есть
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, domain.ErrCategoryExists
		}
		return 0, fmt.Errorf("insert category: %w", err)
	}

	return category_id, nil
}

// List читает все категории, отсортированные по имени.
func (r *CategoryRepo) List(ctx context.Context) ([]domain.Category, error) {
	const query = `
		SELECT id, name
		FROM categories
		ORDER BY name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select categories: %w", err)
	}

	// CollectRows сам закрывает rows и возвращает пустой срез, а не nil, если строк нет
	categories, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Category, error) {
		var c domain.Category
		err := row.Scan(&c.ID, &c.Name)
		return c, err
	})
	if err != nil {
		return nil, fmt.Errorf("scan categories: %w", err)
	}

	return categories, nil
}

// Delete удаляет категорию по id.
// Если строки нет, возвращает domain.ErrCategoryNotFound, если в ней есть карточки, domain.ErrCategoryHasCards.
func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM categories
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		// 23503 foreign_key_violation: на категорию ссылаются карточки (cards.category_id)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return domain.ErrCategoryHasCards
		}
		return fmt.Errorf("delete category: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}
