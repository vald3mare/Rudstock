package category

import (
	"context"
	"fmt"
	"strings"

	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// Update проверяет только переданные поля, теми же правилами, что и Create,
// и отдаёт patch в репозиторий.
func (s *categoryService) Update(ctx context.Context, id int64, in UpdateCategoryInput) (models.Category, error) {
	// пустой PATCH почти всегда ошибка клиента, как и у карточки
	if in.Name == nil {
		return models.Category{}, models.ErrEmptyPatch
	}

	patch := models.CategoryPatch{}

	if in.Name != nil {
		// та же нормализация, что в Create: иначе через PATCH можно было бы
		// завести "Шины " рядом с "Шины"
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return models.Category{}, models.ErrInvalidCategoryName
		}
		patch.Name = &name
	}

	category, err := s.categoryRepo.Update(ctx, id, patch)
	if err != nil {
		return models.Category{}, fmt.Errorf("update category: %w", err)
	}

	return category, nil
}
