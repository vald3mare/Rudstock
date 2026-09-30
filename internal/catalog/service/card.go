package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/internal/catalog/domain"
)

/* РУЧКА (admin)
POST /admin/cards - создание карточки

json {
    category_id: 1,
    description: "Шина летняя 205/55 R16",
    price: 14999,                         // копейки (149.99 ₽)
    photo_url: "s3://photos/photo_1.jpg"  // minio
} -> 201 {card_id: uuid} / error

GET /cards (public) - список карточек
query: category_id?, page?, limit?
-> 200 {items: [...], total}
*/

// Пагинация: 0 означает значение по умолчанию, limit больше максимума срезается до него
const (
	defaultCardsLimit = 20
	maxCardsLimit     = 100
)

// идем сверху вниз
// данные(CardInput) -> CardService.Create(input) -> cardService.Create(input) -> CardRepo.Create() -> postgres.CardRepo.Create()

// CreateCardInput входные данные для создания карточки.
// Контракт: category_id, description, price, photo_url.
type CreateCardInput struct {
	CategoryID  int64  // обязательно, категория должна существовать
	Description string // необязательно
	Price       int64  // обязательно, > 0, копейки
	PhotoURL    string // необязательно, ссылка на объект в minio
}

// ListCardsInput параметры выборки списка карточек.
type ListCardsInput struct {
	CategoryID int64 // необязательно, 0 означает все категории
	Page       int   // необязательно, с 1; 0 означает первую страницу
	Limit      int   // необязательно, 0 означает defaultCardsLimit, максимум maxCardsLimit
}

// CardService бизнес-операции над карточками.
type CardService interface {
	// Create создаёт карточку и возвращает её id.
	// Ошибки: domain.ErrInvalidCategoryID, domain.ErrInvalidPrice, domain.ErrCategoryNotFound.
	Create(ctx context.Context, in CreateCardInput) (uuid.UUID, error)

	// List возвращает страницу карточек и общее число карточек под фильтром.
	// Ошибки: domain.ErrInvalidCategoryID, domain.ErrInvalidPagination.
	List(ctx context.Context, in ListCardsInput) ([]domain.Card, int64, error)
}

// CardRepo что сервису нужно от хранилища (тут описываем требование)
type CardRepo interface {
	// Create сохраняет карточку и возвращает сгенерированный id.
	// Если категории нет, возвращает domain.ErrCategoryNotFound.
	Create(ctx context.Context, card domain.Card) (uuid.UUID, error)

	// List возвращает карточки под фильтром (новые первыми) и их общее число без учёта limit/offset.
	List(ctx context.Context, filter domain.CardFilter) ([]domain.Card, int64, error)
}

// закрытый экземпляр с начинкой
type cardService struct { // делаем сервис неэкспортируемым, чтобы никто не мог создать его копию в обход конструктора
	cardRepo CardRepo // кард репо тоже прячем, тк это внутренности и никто не должен знать как они работают
}

// NewCardService конструктор сервиса (репа не по ссылке, тк интерфейс итак ссылочный тип)
func NewCardService(repo CardRepo) CardService {
	return &cardService{cardRepo: repo}
}

// Create валидирует вход, собирает domain.Card и сохраняет её через репозиторий.
func (s *cardService) Create(ctx context.Context, in CreateCardInput) (uuid.UUID, error) {
	// category_id обязателен: без этой проверки 0 дошёл бы до БД и вернулся как 404
	if in.CategoryID <= 0 {
		return uuid.Nil, domain.ErrInvalidCategoryID
	}

	// цена в копейках, бесплатных и отрицательных карточек не бывает
	if in.Price <= 0 {
		return uuid.Nil, domain.ErrInvalidPrice
	}

	// собираем доменную карточку из инпута, ID и CreatedAt проставляет хранилище
	card := domain.Card{
		CategoryID:  in.CategoryID,
		Description: in.Description,
		Price:       in.Price,
		PhotoURL:    in.PhotoURL,
	}

	card_id, err := s.cardRepo.Create(ctx, card)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create card: %w", err)
	}

	return card_id, nil
}

// List проверяет параметры, переводит page/limit в limit/offset и читает страницу из репозитория.
func (s *cardService) List(ctx context.Context, in ListCardsInput) ([]domain.Card, int64, error) {
	if in.CategoryID < 0 {
		return nil, 0, domain.ErrInvalidCategoryID
	}
	if in.Page < 0 || in.Limit < 0 {
		return nil, 0, domain.ErrInvalidPagination
	}

	page := in.Page
	if page == 0 {
		page = 1
	}

	limit := in.Limit
	if limit == 0 {
		limit = defaultCardsLimit
	}
	limit = min(limit, maxCardsLimit) // защита от ?limit=1000000

	filter := domain.CardFilter{
		CategoryID: in.CategoryID,
		Limit:      limit,
		Offset:     (page - 1) * limit,
	}

	cards, total, err := s.cardRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list cards: %w", err)
	}

	return cards, total, nil
}
