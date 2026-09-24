package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		encodErr := json.NewEncoder(w).Encode(data)
		if encodErr != nil {
			slog.Error("failed to encode response",
				"error", encodErr,
				"method", r.Method,
				"path", r.URL.Path,
			)
		}
	}
}
