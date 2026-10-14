package card

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// Create валидирует вход, собирает models.Card и сохраняет её через репозиторий.
func (s *cardService) Create(ctx context.Context, in CreateCardInput) (uuid.UUID, error) {
	// category_id обязателен: без этой проверки 0 дошёл бы до БД и вернулся как 404
	if in.CategoryID <= 0 {
		return uuid.Nil, models.ErrInvalidCategoryID
	}

	// обрезаем пробелы, как у имени категории: "   " не должно пройти как название
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return uuid.Nil, models.ErrInvalidCardTitle
	}

	// цена в копейках, бесплатных и отрицательных карточек не бывает
	if in.Price <= 0 {
		return uuid.Nil, models.ErrInvalidPrice
	}

	// собираем доменную карточку из инпута, ID и CreatedAt проставляет хранилище
	card := models.Card{
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
