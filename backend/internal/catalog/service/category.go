package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/vald3mare/Rudstock/backend/internal/catalog/domain"
)

/* РУЧКИ
POST /admin/categories (admin) - создание категории
json {
	name: "Шины"
} -> 201 {category_id: int64} / error

GET /categories (public) - список категорий
-> 200 [{id, name}]

DELETE /admin/categories/{id} (admin) - удаление пустой категории
-> 204 без тела / error
*/

// CreateCategoryInput входные данные для создания категории.
type CreateCategoryInput struct {
	Name string // обязательно, не пустое, уникальное
}

// CategoryService бизнес-операции над категориями.
type CategoryService interface {
	// Create создаёт категорию и возвращает её id.
	// Ошибки: domain.ErrInvalidCategoryName, domain.ErrCategoryExists.
	Create(ctx context.Context, in CreateCategoryInput) (int64, error)

	// List возвращает все категории, отсортированные по имени.
	List(ctx context.Context) ([]domain.Category, error)

	// Delete удаляет категорию, если в ней нет карточек.
	// Ошибки: domain.ErrCategoryNotFound, domain.ErrCategoryHasCards.
	Delete(ctx context.Context, id int64) error
}

// CategoryRepo что сервису нужно от хранилища.
type CategoryRepo interface {
	// Create сохраняет категорию и возвращает сгенерированный id.
	// Если имя занято, возвращает domain.ErrCategoryExists.
	Create(ctx context.Context, category domain.Category) (int64, error)

	// List возвращает все категории, отсортированные по имени.
	List(ctx context.Context) ([]domain.Category, error)

	// Delete удаляет категорию.
	// Ошибки: domain.ErrCategoryNotFound, domain.ErrCategoryHasCards.
	Delete(ctx context.Context, id int64) error
}

type categoryService struct {
	categoryRepo CategoryRepo
}

// NewCategoryService конструктор сервиса категорий.
func NewCategoryService(repo CategoryRepo) CategoryService {
	return &categoryService{
		categoryRepo: repo,
	}
}

// Create валидирует имя и сохраняет категорию через репозиторий.
func (s *categoryService) Create(ctx context.Context, in CreateCategoryInput) (int64, error) {
	// обрезаем пробелы: иначе "   " прошло бы проверку,
	// а "Шины" и "Шины " стали бы разными категориями
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return 0, domain.ErrInvalidCategoryName
	}

	category := domain.Category{
		Name: name,
	}

	category_id, err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return 0, fmt.Errorf("create category: %w", err)
	}

	return category_id, nil
}

// List отдаёт категории без фильтров и пагинации: их единицы, витрине нужны все сразу.
func (s *categoryService) List(ctx context.Context) ([]domain.Category, error) {
	categories, err := s.categoryRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	return categories, nil
}

// Delete передаёт удаление в репозиторий. Проверку "в категории есть карточки"
// делает база внешним ключом: отдельный SELECT count перед удалением
// мог бы устареть, пока между ними кто-то создаёт карточку.
func (s *categoryService) Delete(ctx context.Context, id int64) error {
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	return nil
}
