package card

import (
	"net/http"

	"github.com/vald3mare/Rudstock/catalog/internal/platform/httpx"
)

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
