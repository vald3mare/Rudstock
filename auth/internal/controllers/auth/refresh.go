package auth

import (
	"net/http"

	"github.com/vald3mare/Rudstock/auth/internal/platform/httpx"
)

// POST /auth/refresh
func (h *SessionHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	tokens, err := h.svc.Refresh(r.Context(), req.Refresh)
	if err != nil {
		httpx.WriteError(w, r, mapAuthError(err))
		return
	}

	httpx.WriteJSON(w, r, http.StatusOK, toTokensResponse(tokens))
}
