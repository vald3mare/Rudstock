package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/backend/internal/catalog/domain"
)

/* РУЧКА (admin)
POST /admin/cards - создание карточки

json {
    category_id: 1,
    title: "Michelin Primacy 4 205/55 R16",
    description: "Шина летняя 205/55 R16",
    price: 14999,                         // копейки (149.99 ₽)
    photo_url: "s3://photos/photo_1.jpg"  // minio
} -> 201 {card_id: uuid} / error

GET /cards (public) - список карточек
query: category_id?, page?, limit?
-> 200 {items: [...], total}

GET /cards/{id} (public) - одна карточка
-> 200 {id, category_id, title, description, price, photo_url, created_at} / error

PATCH /admin/cards/{id} (admin) - частичное обновление карточки
json {
    price: 12999   // любое подмножество полей из POST, хотя бы одно
} -> 200 карточка целиком / error

DELETE /admin/cards/{id} (admin) - удаление карточки
-> 204 без тела / error
*/

// Пагинация: 0 означает значение по умолчанию, limit больше максимума срезается до него
const (
	defaultCardsLimit = 20
	maxCardsLimit     = 100
)

// идем сверху вниз
// данные(CardInput) -> CardService.Create(input) -> cardService.Create(input) -> CardRepo.Create() -> postgres.CardRepo.Create()

// CreateCardInput входные данные для создания карточки.
// Контракт: category_id, title, description, price, photo_url.
type CreateCardInput struct {
	CategoryID  int64  // обязательно, категория должна существовать
	Title       string // обязательно, не пустое (пробелы по краям обрезаются)
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

// UpdateCardInput входные данные для частичного обновления карточки.
// nil означает "не менять поле"; хотя бы одно поле должно быть задано.
type UpdateCardInput struct {
	CategoryID  *int64  // если задано: > 0, категория должна существовать
	Title       *string // если задано: не пустое (пробелы по краям обрезаются)
	Description *string // если задано: любое, в том числе пустое
	Price       *int64  // если задано: > 0, копейки
	PhotoURL    *string // если задано: любое, в том числе пустое
}

// CardService бизнес-операции над карточками.
type CardService interface {
	// Create создаёт карточку и возвращает её id.
	// Ошибки: domain.ErrInvalidCategoryID, domain.ErrInvalidCardTitle, domain.ErrInvalidPrice, domain.ErrCategoryNotFound.
	Create(ctx context.Context, in CreateCardInput) (uuid.UUID, error)

	// List возвращает страницу карточек и общее число карточек под фильтром.
	// Ошибки: domain.ErrInvalidCategoryID, domain.ErrInvalidPagination.
	List(ctx context.Context, in ListCardsInput) ([]domain.Card, int64, error)

	// Get возвращает карточку по её идентификатору.
	// Ошибки: domain.ErrCardNotFound.
	Get(ctx context.Context, id uuid.UUID) (domain.Card, error)

	// Update меняет переданные поля карточки и возвращает её целиком после изменения.
	// Ошибки: domain.ErrEmptyPatch, domain.ErrInvalidCategoryID, domain.ErrInvalidCardTitle,
	// domain.ErrInvalidPrice, domain.ErrCardNotFound, domain.ErrCategoryNotFound.
	Update(ctx context.Context, id uuid.UUID, in UpdateCardInput) (domain.Card, error)

	// Delete удаляет карточку.
	// Ошибки: domain.ErrCardNotFound.
	Delete(ctx context.Context, id uuid.UUID) error
}

// CardRepo что сервису нужно от хранилища (тут описываем требование)
type CardRepo interface {
	// Create сохраняет карточку и возвращает сгенерированный id.
	// Если категории нет, возвращает domain.ErrCategoryNotFound.
	Create(ctx context.Context, card domain.Card) (uuid.UUID, error)

	// List возвращает карточки под фильтром (новые первыми) и их общее число без учёта limit/offset.
	List(ctx context.Context, filter domain.CardFilter) ([]domain.Card, int64, error)

	// Get возвращает карточку по её идентификатору.
	// Ошибки: domain.ErrCardNotFound.
	Get(ctx context.Context, id uuid.UUID) (domain.Card, error)

	// Update применяет patch и возвращает карточку после изменения.
	// Ошибки: domain.ErrCardNotFound, domain.ErrCategoryNotFound.
	Update(ctx context.Context, id uuid.UUID, patch domain.CardPatch) (domain.Card, error)

	// Delete удаляет карточку.
	// Ошибки: domain.ErrCardNotFound.
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

// Create валидирует вход, собирает domain.Card и сохраняет её через репозиторий.
func (s *cardService) Create(ctx context.Context, in CreateCardInput) (uuid.UUID, error) {
	// category_id обязателен: без этой проверки 0 дошёл бы до БД и вернулся как 404
	if in.CategoryID <= 0 {
		return uuid.Nil, domain.ErrInvalidCategoryID
	}

	// обрезаем пробелы, как у имени категории: "   " не должно пройти как название
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return uuid.Nil, domain.ErrInvalidCardTitle
	}

	// цена в копейках, бесплатных и отрицательных карточек не бывает
	if in.Price <= 0 {
		return uuid.Nil, domain.ErrInvalidPrice
	}

	// собираем доменную карточку из инпута, ID и CreatedAt проставляет хранилище
	card := domain.Card{
		CategoryID:  in.CategoryID,
		Title:       title,
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

// Get передаёт поиск в репозиторий. Отдельной проверки на uuid.Nil нет:
// база его не найдёт, и клиент получит 404, как и для любого несуществующего id.
func (s *cardService) Get(ctx context.Context, id uuid.UUID) (domain.Card, error) {
	card, err := s.cardRepo.Get(ctx, id)
	if err != nil {
		return domain.Card{}, fmt.Errorf("get card: %w", err)
	}

	return card, nil
}

// Update проверяет только переданные поля, теми же правилами, что и Create,
// и отдаёт patch в репозиторий.
func (s *cardService) Update(ctx context.Context, id uuid.UUID, in UpdateCardInput) (domain.Card, error) {
	// пустой PATCH почти всегда ошибка клиента (опечатка в имени поля и т.п.),
	// молча вернуть 200 без изменений значило бы её спрятать
	if in.CategoryID == nil && in.Title == nil && in.Description == nil && in.Price == nil && in.PhotoURL == nil {
		return domain.Card{}, domain.ErrEmptyPatch
	}

	if in.CategoryID != nil && *in.CategoryID <= 0 {
		return domain.Card{}, domain.ErrInvalidCategoryID
	}

	// title, если передан, нормализуем так же, как в Create, и кладём в patch уже обрезанным
	var title *string
	if in.Title != nil {
		t := strings.TrimSpace(*in.Title)
		if t == "" {
			return domain.Card{}, domain.ErrInvalidCardTitle
		}
		title = &t
	}

	if in.Price != nil && *in.Price <= 0 {
		return domain.Card{}, domain.ErrInvalidPrice
	}

	patch := domain.CardPatch{
		CategoryID:  in.CategoryID,
		Title:       title,
		Description: in.Description,
		Price:       in.Price,
		PhotoURL:    in.PhotoURL,
	}

	card, err := s.cardRepo.Update(ctx, id, patch)
	if err != nil {
		return domain.Card{}, fmt.Errorf("update card: %w", err)
	}

	return card, nil
}

// Delete передаёт удаление в репозиторий.
// Пока удаляем физически: заказов ещё нет, архивировать (status) будем, когда появятся.
func (s *cardService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.cardRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete card: %w", err)
	}

	return nil
}
