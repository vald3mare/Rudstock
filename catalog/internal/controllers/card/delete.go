package card

import (
	"net/http"

	"github.com/vald3mare/Rudstock/catalog/internal/platform/httpx"
)

// DELETE /admin/cards/{id}
func (h *AdminCardHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
