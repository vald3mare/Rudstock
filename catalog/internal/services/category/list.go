package category

import (
	"context"
	"fmt"

	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// List отдаёт категории без фильтров и пагинации: их единицы, витрине нужны все сразу.
func (s *categoryService) List(ctx context.Context) ([]models.Category, error) {
	categories, err := s.categoryRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	return categories, nil
}
