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
*/

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

// CardService бизнес-операции над карточками.
type CardService interface {
	// Create создаёт карточку и возвращает её id.
	// Ошибки: domain.ErrInvalidCategoryID, domain.ErrInvalidPrice, domain.ErrCategoryNotFound.
	Create(ctx context.Context, in CreateCardInput) (uuid.UUID, error)
}

// CardRepo что сервису нужно от хранилища (тут описываем требование)
type CardRepo interface {
	// Create сохраняет карточку и возвращает сгенерированный id.
	// Если категории нет, возвращает domain.ErrCategoryNotFound.
	Create(ctx context.Context, card domain.Card) (uuid.UUID, error)
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
