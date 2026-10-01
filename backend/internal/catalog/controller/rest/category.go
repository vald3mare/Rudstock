package rest

import (
	"errors"
	"net/http"

	"github.com/vald3mare/Rudstock/backend/internal/catalog/domain"
	"github.com/vald3mare/Rudstock/backend/internal/catalog/service"
	"github.com/vald3mare/Rudstock/backend/internal/platform/errs"
	"github.com/vald3mare/Rudstock/backend/internal/platform/httpx"
)

// CategoryHandler HTTP-ручки категорий: разбирает запрос, зовёт сервис, пишет ответ.
type CategoryHandler struct {
	svc service.CategoryService
}

// NewCategoryHandler конструктор хендлера категорий.
func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: s}
}

// POST /admin/categories
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	id, err := h.svc.Create(r.Context(), req.toInput())
	if err != nil {
		httpx.WriteError(w, r, mapCategoryError(err))
		return
	}

	httpx.WriteJSON(w, r, http.StatusCreated, CreateCategoryResponse{CategoryID: id})
}

// GET /categories
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.List(r.Context())
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, toCategoryResponses(categories))
}

// mapCategoryError переводит доменные ошибки в errs, чтобы httpx выбрал правильный код.
// Всё неизвестное уходит как есть, WriteError ответит 500.
func mapCategoryError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidCategoryName):
		return errs.InvalidInput("invalid-category-name", "Category name must not be empty", err)
	case errors.Is(err, domain.ErrCategoryExists):
		return errs.Conflict("category-exists", "Category with this name already exists", err)
	default:
		return err
	}
}
