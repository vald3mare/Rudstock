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

// cardColumns и scanCard описывают карточку один раз для всех SELECT/RETURNING:
// новое поле добавляется в двух местах рядом, а не в каждом запросе,
// и порядок колонок не разойдётся с порядком Scan.
const cardColumns = `id, category_id, title, description, price, photo_url, created_at`

// scanCard читает строку с колонками cardColumns. pgx.Row подходит и для QueryRow, и для строк CollectRows.
func scanCard(row pgx.Row) (domain.Card, error) {
	var c domain.Card
	err := row.Scan(&c.ID, &c.CategoryID, &c.Title, &c.Description, &c.Price, &c.PhotoURL, &c.CreatedAt)
	return c, err
}

// вот и сама реализация создания карточки, именно от сюда мы и ходим в бд
func (r *CardRepo) Create(ctx context.Context, card domain.Card) (uuid.UUID, error) {
	// id и created_at не передаем, их проставляют DEFAULT'ы в таблице
	const query = `
		INSERT INTO cards (category_id, title, description, price, photo_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	var card_id uuid.UUID
	err := r.pool.QueryRow(ctx, query,
		card.CategoryID,
		card.Title,
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
		SELECT ` + cardColumns + `
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
		return scanCard(row)
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

// Get читает карточку по id. Если строки нет, возвращает domain.ErrCardNotFound.
func (r *CardRepo) Get(ctx context.Context, id uuid.UUID) (domain.Card, error) {
	const query = `
		SELECT ` + cardColumns + `
		FROM cards
		WHERE id = $1`

	c, err := scanCard(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Card{}, domain.ErrCardNotFound
		}
		return domain.Card{}, fmt.Errorf("select card: %w", err)
	}

	return c, nil
}

// Update применяет patch одним запросом и возвращает карточку после изменения.
// Если карточки нет, возвращает domain.ErrCardNotFound, если новой категории нет, domain.ErrCategoryNotFound.
func (r *CardRepo) Update(ctx context.Context, id uuid.UUID, patch domain.CardPatch) (domain.Card, error) {
	// nil-указатель уходит в базу как NULL, а COALESCE(NULL, колонка) оставляет старое значение.
	// Так не нужно собирать SET строкой под каждый набор полей.
	// RETURNING отдаёт карточку целиком, второй SELECT не нужен.
	const query = `
		UPDATE cards
		SET category_id = COALESCE($2, category_id),
		    title       = COALESCE($3, title),
		    description = COALESCE($4, description),
		    price       = COALESCE($5, price),
		    photo_url   = COALESCE($6, photo_url)
		WHERE id = $1
		RETURNING ` + cardColumns

	c, err := scanCard(r.pool.QueryRow(ctx, query,
		id,
		patch.CategoryID,
		patch.Title,
		patch.Description,
		patch.Price,
		patch.PhotoURL,
	))
	if err != nil {
		// UPDATE без подходящей строки ничего не возвращает: карточки с таким id нет
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Card{}, domain.ErrCardNotFound
		}
		// 23503 foreign_key_violation: новой категории с таким id нет
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return domain.Card{}, domain.ErrCategoryNotFound
		}
		return domain.Card{}, fmt.Errorf("update card: %w", err)
	}

	return c, nil
}

// Delete удаляет карточку по id. Если строки нет, возвращает domain.ErrCardNotFound.
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
		return domain.ErrCardNotFound
	}

	return nil
}
