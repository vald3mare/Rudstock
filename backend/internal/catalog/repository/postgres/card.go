package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vald3mare/Rudstock/backend/internal/catalog/domain"
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

// List читает страницу карточек и отдельным запросом общее число под фильтром.
func (r *CardRepo) List(ctx context.Context, filter domain.CardFilter) ([]domain.Card, int64, error) {
	// $1 = 0 выключает фильтр по категории, так не нужно собирать SQL строкой.
	// id во втором ключе сортировки: у карточек с одинаковым created_at
	// порядок должен быть стабильным, иначе они будут прыгать между страницами.
	const listQuery = `
		SELECT id, category_id, description, price, photo_url, created_at
		FROM cards
		WHERE ($1::bigint = 0 OR category_id = $1)
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3`

	const countQuery = `
		SELECT count(*)
		FROM cards
		WHERE ($1::bigint = 0 OR category_id = $1)`

	rows, err := r.pool.Query(ctx, listQuery, filter.CategoryID, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("select cards: %w", err)
	}

	// CollectRows сам закрывает rows и возвращает пустой срез, а не nil, если строк нет
	cards, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Card, error) {
		var c domain.Card
		err := row.Scan(&c.ID, &c.CategoryID, &c.Description, &c.Price, &c.PhotoURL, &c.CreatedAt)
		return c, err
	})
	if err != nil {
		return nil, 0, fmt.Errorf("scan cards: %w", err)
	}

	// total считаем отдельно: count(*) OVER () в первом запросе вернул бы 0,
	// если страница за пределами списка и строк нет
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, filter.CategoryID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count cards: %w", err)
	}

	return cards, total, nil
}
