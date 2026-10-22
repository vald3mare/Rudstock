package auth

import (
	"net/http"

	"github.com/vald3mare/Rudstock/auth/internal/platform/httpx"
)

// POST /auth/register
func (h *CredentialsHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	res, err := h.svc.Register(r.Context(), req.toInput())
	if err != nil {
		httpx.WriteError(w, r, mapAuthError(err))
		return
	}

	httpx.WriteJSON(w, r, http.StatusCreated, toAuthResponse(res))
}
