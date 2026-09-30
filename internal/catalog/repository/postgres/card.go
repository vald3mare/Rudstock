package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vald3mare/Rudstock/internal/catalog/domain"
)

// репа может быть открытой, тк сервис ничего не знает о репозитории, единственное где мы дергаем его, это в неэкспортируемом поле при вызове конструктора
type CardRepo struct {
	pool *pgxpool.Pool
}

// NewCardRepo конструктор репозитория карточек.
func NewCardRepo(pool *pgxpool.Pool) *CardRepo {
	return &CardRepo{
		pool: pool,
	}
}

// вот и сама реализация создания карточки, именно от сюда мы и ходим в бд
func (r *CardRepo) Create(ctx context.Context, card domain.Card) (uuid.UUID, error) {
	// id и created_at не передаем, их проставляют DEFAULT'ы в таблице
	const query = `
		INSERT INTO cards (category_id, description, price, photo_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var card_id uuid.UUID
	err := r.pool.QueryRow(ctx, query,
		card.CategoryID,
		card.Description,
		card.Price,
		card.PhotoURL,
	).Scan(&card_id)
	if err != nil {
		// 23503 foreign_key_violation: категории с таким id нет
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return uuid.Nil, domain.ErrCategoryNotFound
		}
		return uuid.Nil, fmt.Errorf("insert card: %w", err)
	}

	return card_id, nil
}
