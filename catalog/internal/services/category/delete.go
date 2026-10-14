package category

import (
	"context"
	"fmt"
)

// Delete передаёт удаление в репозиторий. Проверку "в категории есть карточки"
// делает база внешним ключом: отдельный SELECT count перед удалением
// мог бы устареть, пока между ними кто-то создаёт карточку.
func (s *categoryService) Delete(ctx context.Context, id int64) error {
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	return nil
}
