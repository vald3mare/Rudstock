package card

import (
	"context"
	"fmt"

	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// List проверяет параметры, переводит page/limit в limit/offset и читает страницу из репозитория.
func (s *cardService) List(ctx context.Context, in ListCardsInput) ([]models.Card, int64, error) {
	if in.CategoryID < 0 {
		return nil, 0, models.ErrInvalidCategoryID
	}
	if in.Page < 0 || in.Limit < 0 {
		return nil, 0, models.ErrInvalidPagination
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

	filter := models.CardFilter{
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
