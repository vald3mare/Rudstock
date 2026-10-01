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

PATCH /admin/categories/{id} (admin) - переименование категории
json {
	name: "Летние шины"
} -> 200 {id, name} / error

DELETE /admin/categories/{id} (admin) - удаление пустой категории
-> 204 без тела / error
*/

// CreateCategoryInput входные данные для создания категории.
type CreateCategoryInput struct {
	Name string // обязательно, не пустое, уникальное
}

// UpdateCategoryInput входные данные для частичного обновления категории.
// nil означает "не менять поле"; хотя бы одно поле должно быть задано.
// Поле пока одно, но указатель оставляем как у карточки: новые поля (sort и т.п.) добавятся без смены подхода.
type UpdateCategoryInput struct {
	Name *string // если задано: не пустое, уникальное
}

// CategoryService бизнес-операции над категориями.
type CategoryService interface {
	// Create создаёт категорию и возвращает её id.
	// Ошибки: domain.ErrInvalidCategoryName, domain.ErrCategoryExists.
	Create(ctx context.Context, in CreateCategoryInput) (int64, error)

	// List возвращает все категории, отсортированные по имени.
	List(ctx context.Context) ([]domain.Category, error)

	// Update меняет переданные поля категории и возвращает её после изменения.
	// Ошибки: domain.ErrEmptyPatch, domain.ErrInvalidCategoryName,
	// domain.ErrCategoryNotFound, domain.ErrCategoryExists.
	Update(ctx context.Context, id int64, in UpdateCategoryInput) (domain.Category, error)

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

	// Update применяет patch и возвращает категорию после изменения.
	// Ошибки: domain.ErrCategoryNotFound, domain.ErrCategoryExists.
	Update(ctx context.Context, id int64, patch domain.CategoryPatch) (domain.Category, error)

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

// Update проверяет только переданные поля, теми же правилами, что и Create,
// и отдаёт patch в репозиторий.
func (s *categoryService) Update(ctx context.Context, id int64, in UpdateCategoryInput) (domain.Category, error) {
	// пустой PATCH почти всегда ошибка клиента, как и у карточки
	if in.Name == nil {
		return domain.Category{}, domain.ErrEmptyPatch
	}

	patch := domain.CategoryPatch{}

	if in.Name != nil {
		// та же нормализация, что в Create: иначе через PATCH можно было бы
		// завести "Шины " рядом с "Шины"
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return domain.Category{}, domain.ErrInvalidCategoryName
		}
		patch.Name = &name
	}

	category, err := s.categoryRepo.Update(ctx, id, patch)
	if err != nil {
		return domain.Category{}, fmt.Errorf("update category: %w", err)
	}

	return category, nil
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
