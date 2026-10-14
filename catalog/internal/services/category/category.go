package category

import (
	"context"

	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

// CategoryService бизнес-операции над категориями.
type CategoryService interface {
	// Create создаёт категорию и возвращает её id.s
	// Ошибки: models.ErrInvalidCategoryName, models.ErrCategoryExists.
	Create(ctx context.Context, in CreateCategoryInput) (int64, error)

	// List возвращает все категории, отсортированные по имени.
	List(ctx context.Context) ([]models.Category, error)

	// Update меняет переданные поля категории и возвращает её после изменения.
	// Ошибки: models.ErrEmptyPatch, models.ErrInvalidCategoryName,
	// models.ErrCategoryNotFound, models.ErrCategoryExists.
	Update(ctx context.Context, id int64, in UpdateCategoryInput) (models.Category, error)

	// Delete удаляет категорию, если в ней нет карточек.
	// Ошибки: models.ErrCategoryNotFound, models.ErrCategoryHasCards.
	Delete(ctx context.Context, id int64) error
}

// CategoryRepo что сервису нужно от хранилища.
type CategoryRepo interface {
	// Create сохраняет категорию и возвращает сгенерированный id.
	// Если имя занято, возвращает models.ErrCategoryExists.
	Create(ctx context.Context, category models.Category) (int64, error)

	// List возвращает все категории, отсортированные по имени.
	List(ctx context.Context) ([]models.Category, error)

	// Update применяет patch и возвращает категорию после изменения.
	// Ошибки: models.ErrCategoryNotFound, models.ErrCategoryExists.
	Update(ctx context.Context, id int64, patch models.CategoryPatch) (models.Category, error)

	// Delete удаляет категорию.
	// Ошибки: models.ErrCategoryNotFound, models.ErrCategoryHasCards.
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
