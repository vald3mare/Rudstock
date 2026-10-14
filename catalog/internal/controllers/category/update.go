package category

import (
	"net/http"

	"github.com/vald3mare/Rudstock/catalog/internal/platform/httpx"
)

// PATCH /admin/categories/{id}
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathInt64(r, "id")
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	var req UpdateCategoryRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	category, err := h.svc.Update(r.Context(), id, req.toInput())
	if err != nil {
		httpx.WriteError(w, r, mapCategoryError(err))
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, toCategoryResponse(category))
}
