package card

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Delete передаёт удаление в репозиторий.
// Пока удаляем физически: заказов ещё нет, архивировать (status) будем, когда появятся.
func (s *cardService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.cardRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete card: %w", err)
	}

	return nil
}
