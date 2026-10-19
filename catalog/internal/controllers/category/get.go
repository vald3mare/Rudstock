package category

import (
	"net/http"

	"github.com/vald3mare/Rudstock/catalog/internal/platform/httpx"
)

// GET /categories
func (h *PublicCategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.List(r.Context())
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, toCategoryResponses(categories))
}
