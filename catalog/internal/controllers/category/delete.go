package category

import (
	"net/http"

	"github.com/vald3mare/Rudstock/catalog/internal/platform/httpx"
)

// DELETE /admin/categories/{id}
func (h *AdminCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathInt64(r, "id")
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		httpx.WriteError(w, r, mapCategoryError(err))
		return
	}

	// 204 без тела: удалённую категорию возвращать незачем
	w.WriteHeader(http.StatusNoContent)
}
