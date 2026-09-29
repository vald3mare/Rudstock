package rest

import (
	"errors"
	"net/http"

	"github.com/vald3mare/Rudstock/internal/catalog/domain"
	"github.com/vald3mare/Rudstock/internal/catalog/service"
	"github.com/vald3mare/Rudstock/internal/platform/errs"
	"github.com/vald3mare/Rudstock/internal/platform/httpx"
)

type CardHandler struct {
	svc service.CardService
}

func NewCardHandler(s service.CardService) *CardHandler {
	return &CardHandler{svc: s}
}

// POST /admin/cards
func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCardRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	id, err := h.svc.Create(r.Context(), service.CreateCardInput{
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Price:       req.Price,
		PhotoURL:    req.PhotoURL,
	})
	if err != nil {
		httpx.WriteError(w, r, mapCardError(err))
		return
	}

	httpx.WriteJSON(w, r, http.StatusCreated, CreateCardResponse{CardID: id})
}

// mapCardError переводит доменные ошибки в errs, чтобы httpx выбрал правильный код.
// Всё неизвестное уходит как есть, WriteError ответит 500.
func mapCardError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidPrice):
		return errs.InvalidInput("invalid-price", "Price must be greater than zero", err)
	case errors.Is(err, domain.ErrCategoryNotFound):
		return errs.NotFound("category-not-found", "Category not found", err)
	default:
		return err
	}
}
