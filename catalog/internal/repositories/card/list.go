package card

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

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

// List читает страницу карточек и отдельным запросом общее число под фильтром.
func (r *CardRepo) List(ctx context.Context, filter models.CardFilter) (cards []models.Card, total int64, err error) {
	rows, err := r.pool.Query(ctx, listQuery, filter.CategoryID, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("select cards: %w", err)
	}

	// CollectRows сам закрывает rows и возвращает пустой срез, а не nil, если строк нет
	cards, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Card, error) {
		return scanCard(row)
	})
	if err != nil {
		return nil, 0, fmt.Errorf("scan cards: %w", err)
	}

	// total считаем отдельно: count(*) OVER () в первом запросе вернул бы 0,
	// если страница за пределами списка и строк нет
	if err := r.pool.QueryRow(ctx, countQuery, filter.CategoryID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count cards: %w", err)
	}

	return cards, total, nil
}
