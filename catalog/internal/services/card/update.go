package card

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// Update проверяет только переданные поля, теми же правилами, что и Create,
// и отдаёт patch в репозиторий.
func (s *cardService) Update(ctx context.Context, id uuid.UUID, in UpdateCardInput) (models.Card, error) {
	// пустой PATCH почти всегда ошибка клиента (опечатка в имени поля и т.п.),
	// молча вернуть 200 без изменений значило бы её спрятать
	if in.CategoryID == nil && in.Title == nil && in.Description == nil && in.Price == nil && in.PhotoURL == nil {
		return models.Card{}, models.ErrEmptyPatch
	}

	if in.CategoryID != nil && *in.CategoryID <= 0 {
		return models.Card{}, models.ErrInvalidCategoryID
	}

	// title, если передан, нормализуем так же, как в Create, и кладём в patch уже обрезанным
	var title *string
	if in.Title != nil {
		t := strings.TrimSpace(*in.Title)
		if t == "" {
			return models.Card{}, models.ErrInvalidCardTitle
		}
		title = &t
	}

	if in.Price != nil && *in.Price <= 0 {
		return models.Card{}, models.ErrInvalidPrice
	}

	patch := models.CardPatch{
		CategoryID:  in.CategoryID,
		Title:       title,
		Description: in.Description,
		Price:       in.Price,
		PhotoURL:    in.PhotoURL,
	}

	card, err := s.cardRepo.Update(ctx, id, patch)
	if err != nil {
		return models.Card{}, fmt.Errorf("update card: %w", err)
	}

	return card, nil
}
