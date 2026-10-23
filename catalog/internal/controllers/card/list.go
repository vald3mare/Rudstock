package card

import (
	"net/http"

	"github.com/vald3mare/Rudstock/catalog/internal/platform/httpx"
	"github.com/vald3mare/Rudstock/catalog/internal/services/card"
)

// GET /cards?category_id=&page=&limit=
func (h *PublicCardHandler) List(w http.ResponseWriter, r *http.Request) {
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

	cards, total, err := h.svc.List(r.Context(), card.ListCardsInput{
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
