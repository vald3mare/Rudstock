package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/internal/catalog/domain"
)

/* ПУБЛИЧНЫЕ РУЧКИ
/card/create - создание карточик

POST json {
    category_id: 1,
    description: "",
    price: 100.00,
    photo_url: "s3://photo_1" (minio)
} -> card_id / error
*/

// идем сверху вниз
// данные(CardInput) -> CardService.Create(input) -> cardService.Create(input) -> CardRepo.Create() -> postgres.CardRepo.Create()

// структура входных данных для создания карточки
// контракт: category_id, description, price, photo_url
type CreateCardInput struct {
	CategoryID  int64
	Description string
	Price       int64 //копейки
	PhotoURL    string
}

// общий интерфейс
type CardService interface {
	Create(ctx context.Context, in CreateCardInput) (uuid.UUID, error)
}

// что сервису нужно для хранилища (тут описываем требование)
type CardRepo interface {
	Create(ctx context.Context, card domain.Card) (uuid.UUID, error)
}

// закрытый экземпляр с начинкой
type cardService struct { // делаем сервис неэкспортируемым, чтобы никто не мог создать его копию в обход конструктора
	cardRepo CardRepo // кард репо тоже прячем, тк это внутренности и никто не должен знать как они работают
}

// конструктор чтобы создавать сервис (репа не по ссылке, тк итерфейс итак ссылочный тип)
func NewCardService(repo CardRepo) CardService {
	return &cardService{cardRepo: repo}
}

// создаем карточку
func (s *cardService) Create(ctx context.Context, in CreateCardInput) (uuid.UUID, error) {
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

	id, err := s.cardRepo.Create(ctx, card)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create card: %w", err)
	}

	return id, nil
}
