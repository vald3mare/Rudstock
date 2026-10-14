package card

import (
	"errors"

	"github.com/vald3mare/Rudstock/catalog/internal/models"
	"github.com/vald3mare/Rudstock/catalog/internal/platform/errs"
)

// mapCardError переводит доменные ошибки в errs, чтобы httpx выбрал правильный код.
// Всё неизвестное уходит как есть, WriteError ответит 500.
func mapCardError(err error) error {
	switch {
	case errors.Is(err, models.ErrInvalidPrice):
		return errs.InvalidInput("invalid-price", "Price must be greater than zero", err)
	case errors.Is(err, models.ErrInvalidCardTitle):
		return errs.InvalidInput("invalid-title", "Title must not be empty", err)
	case errors.Is(err, models.ErrInvalidCategoryID):
		return errs.InvalidInput("invalid-category-id", "category_id must be a positive integer", err)
	case errors.Is(err, models.ErrInvalidPagination):
		return errs.InvalidInput("invalid-pagination", "page and limit must not be negative", err)
	case errors.Is(err, models.ErrEmptyPatch):
		return errs.InvalidInput("empty-patch", "At least one field must be provided", err)
	case errors.Is(err, models.ErrCategoryNotFound):
		return errs.NotFound("category-not-found", "Category not found", err)
	case errors.Is(err, models.ErrCardNotFound):
		return errs.NotFound("card-not-found", "Card not found", err)
	default:
		return err
	}
}
