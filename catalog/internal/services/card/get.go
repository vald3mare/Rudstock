package card

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// Get передаёт поиск в репозиторий. Отдельной проверки на uuid.Nil нет:
// база его не найдёт, и клиент получит 404, как и для любого несуществующего id.
func (s *cardService) Get(ctx context.Context, id uuid.UUID) (models.Card, error) {
	card, err := s.cardRepo.Get(ctx, id)
	if err != nil {
		return models.Card{}, fmt.Errorf("get card: %w", err)
	}

	return card, nil
}
