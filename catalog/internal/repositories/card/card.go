package card

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// CardRepo что умеет репозиторий карточек; конструктор отдаёт его, а не структуру.
// Требование сервиса описано отдельно, в services/card.CardRepo.
type CardRepo interface {
	// Create сохраняет карточку и возвращает сгенерированный id.
	// Ошибки: models.ErrCategoryNotFound.
	Create(ctx context.Context, card models.Card) (uuid.UUID, error)

	// List возвращает карточки под фильтром (новые первыми) и их общее число без учёта limit/offset.
	List(ctx context.Context, filter models.CardFilter) (cards []models.Card, total int64, err error)

	// Get возвращает карточку по её идентификатору.
	// Ошибки: models.ErrCardNotFound.
	Get(ctx context.Context, id uuid.UUID) (models.Card, error)

	// Update применяет patch и возвращает карточку после изменения.
	// Ошибки: models.ErrCardNotFound, models.ErrCategoryNotFound.
	Update(ctx context.Context, id uuid.UUID, patch models.CardPatch) (models.Card, error)

	// Delete удаляет карточку.
	// Ошибки: models.ErrCardNotFound.
	Delete(ctx context.Context, id uuid.UUID) error
}

// cardRepo закрытая: создаётся только через NewCardRepo, наружу отдаётся интерфейс CardRepo, пул снаружи не виден
type cardRepo struct {
	pool *pgxpool.Pool
}

// NewCardRepo конструктор репозитория карточек.
func NewCardRepo(pool *pgxpool.Pool) CardRepo {
	return &cardRepo{
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
