package category

import (
	"context"
	"fmt"
	"strings"

	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// Create валидирует имя и сохраняет категорию через репозиторий.
func (s *categoryService) Create(ctx context.Context, in CreateCategoryInput) (int64, error) {
	// обрезаем пробелы: иначе "   " прошло бы проверку,
	// а "Шины" и "Шины " стали бы разными категориями
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return 0, models.ErrInvalidCategoryName
	}

	category := models.Category{
		Name: name,
	}

	category_id, err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return 0, fmt.Errorf("create category: %w", err)
	}

	return category_id, nil
}
