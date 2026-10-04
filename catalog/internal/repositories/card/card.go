package card

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
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
func scanCard(row pgx.Row) (models.Card, error) {
	var c models.Card
	err := row.Scan(&c.ID, &c.CategoryID, &c.Title, &c.Description, &c.Price, &c.PhotoURL, &c.CreatedAt)
	return c, err
}
