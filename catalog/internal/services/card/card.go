package card

import (
	"context"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// CardService бизнес-операции над карточками.
type CardService interface {
	// Create создаёт карточку и возвращает её id.
	// Ошибки: models.ErrInvalidCategoryID, models.ErrInvalidCardTitle, models.ErrInvalidPrice, models.ErrCategoryNotFound.
	Create(ctx context.Context, in CreateCardInput) (uuid.UUID, error)

	// List возвращает страницу карточек и общее число карточек под фильтром.
	// Ошибки: models.ErrInvalidCategoryID, models.ErrInvalidPagination.
	List(ctx context.Context, in ListCardsInput) ([]models.Card, int64, error)

	// Get возвращает карточку по её идентификатору.
	// Ошибки: models.ErrCardNotFound.
	Get(ctx context.Context, id uuid.UUID) (models.Card, error)

	// Update меняет переданные поля карточки и возвращает её целиком после изменения.
	// Ошибки: models.ErrEmptyPatch, models.ErrInvalidCategoryID, models.ErrInvalidCardTitle,
	// models.ErrInvalidPrice, models.ErrCardNotFound, models.ErrCategoryNotFound.
	Update(ctx context.Context, id uuid.UUID, in UpdateCardInput) (models.Card, error)

	// Delete удаляет карточку.
	// Ошибки: models.ErrCardNotFound.
	Delete(ctx context.Context, id uuid.UUID) error
}

// CardRepo что сервису нужно от хранилища (тут описываем требование)
type CardRepo interface {
	// Create сохраняет карточку и возвращает сгенерированный id.
	// Если категории нет, возвращает models.ErrCategoryNotFound.
	Create(ctx context.Context, card models.Card) (uuid.UUID, error)

	// List возвращает карточки под фильтром (новые первыми) и их общее число без учёта limit/offset.
	List(ctx context.Context, filter models.CardFilter) ([]models.Card, int64, error)

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

// закрытый экземпляр с начинкой
type cardService struct { // делаем сервис неэкспортируемым, чтобы никто не мог создать его копию в обход конструктора
	cardRepo CardRepo // кард репо тоже прячем, тк это внутренности и никто не должен знать как они работают
}

// NewCardService конструктор сервиса (репа не по ссылке, тк интерфейс итак ссылочный тип)
func NewCardService(repo CardRepo) CardService {
	return &cardService{cardRepo: repo}
}
