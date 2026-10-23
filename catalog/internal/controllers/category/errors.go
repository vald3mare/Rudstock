package category

import (
	"errors"

	"github.com/vald3mare/Rudstock/catalog/internal/models"
	"github.com/vald3mare/Rudstock/catalog/internal/platform/errs"
)

// mapCategoryError переводит доменные ошибки в errs, чтобы httpx выбрал правильный код.
// Всё неизвестное уходит как есть, WriteError ответит 500.
func mapCategoryError(err error) error {
	switch {
	case errors.Is(err, models.ErrInvalidCategoryName):
		return errs.InvalidInput("invalid-category-name", "Category name must not be empty", err)
	case errors.Is(err, models.ErrCategoryExists):
		return errs.Conflict("category-exists", "Category with this name already exists", err)
	case errors.Is(err, models.ErrCategoryHasCards):
		return errs.Conflict("category-has-cards", "Category has cards, move or delete them first", err)
	case errors.Is(err, models.ErrCategoryNotFound):
		return errs.NotFound("category-not-found", "Category not found", err)
	case errors.Is(err, models.ErrEmptyPatch):
		return errs.InvalidInput("empty-patch", "At least one field must be provided", err)
	default:
		return err
	}
}
