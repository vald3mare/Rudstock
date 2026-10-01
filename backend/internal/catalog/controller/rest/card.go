package rest

import (
	"errors"
	"net/http"

	"github.com/vald3mare/Rudstock/backend/internal/catalog/domain"
	"github.com/vald3mare/Rudstock/backend/internal/catalog/service"
	"github.com/vald3mare/Rudstock/backend/internal/platform/errs"
	"github.com/vald3mare/Rudstock/backend/internal/platform/httpx"
)

// CardHandler HTTP-ручки карточек: разбирает запрос, зовёт сервис, пишет ответ.
type CardHandler struct {
	svc service.CardService
}

// NewCardHandler конструктор хендлера карточек.
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

	id, err := h.svc.Create(r.Context(), req.toInput())
	if err != nil {
		httpx.WriteError(w, r, mapCardError(err))
		return
	}

	httpx.WriteJSON(w, r, http.StatusCreated, CreateCardResponse{CardID: id})
}

// GET /cards?category_id=&page=&limit=
func (h *CardHandler) List(w http.ResponseWriter, r *http.Request) {
	categoryID, err := httpx.QueryInt64(r, "category_id")
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	page, err := httpx.QueryInt64(r, "page")
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}
	limit, err := httpx.QueryInt64(r, "limit")
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	cards, total, err := h.svc.List(r.Context(), service.ListCardsInput{
		CategoryID: categoryID,
		Page:       int(page),
		Limit:      int(limit),
	})
	if err != nil {
		httpx.WriteError(w, r, mapCardError(err))
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, ListCardsResponse{
		Items: toCardResponses(cards),
		Total: total,
	})
}

// GET /cards/{id}
func (h *CardHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	card, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, r, mapCardError(err))
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, toCardResponse(card))
}

// PATCH /admin/cards/{id}
func (h *CardHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	var req UpdateCardRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	card, err := h.svc.Update(r.Context(), id, req.toInput())
	if err != nil {
		httpx.WriteError(w, r, mapCardError(err))
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, toCardResponse(card))
}

// DELETE /admin/cards/{id}
func (h *CardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		httpx.WriteError(w, r, mapCardError(err))
		return
	}

	// 204 без тела: удалённую карточку возвращать незачем
	w.WriteHeader(http.StatusNoContent)
}

// mapCardError переводит доменные ошибки в errs, чтобы httpx выбрал правильный код.
// Всё неизвестное уходит как есть, WriteError ответит 500.
func mapCardError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidPrice):
		return errs.InvalidInput("invalid-price", "Price must be greater than zero", err)
	case errors.Is(err, domain.ErrInvalidCategoryID):
		return errs.InvalidInput("invalid-category-id", "category_id must be a positive integer", err)
	case errors.Is(err, domain.ErrInvalidPagination):
		return errs.InvalidInput("invalid-pagination", "page and limit must not be negative", err)
	case errors.Is(err, domain.ErrEmptyPatch):
		return errs.InvalidInput("empty-patch", "At least one field must be provided", err)
	case errors.Is(err, domain.ErrCategoryNotFound):
		return errs.NotFound("category-not-found", "Category not found", err)
	case errors.Is(err, domain.ErrCardNotFound):
		return errs.NotFound("card-not-found", "Card not found", err)
	default:
		return err
	}
}
