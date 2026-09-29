package rest

import "net/http"

func NewRouter(cards *CardHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /admin/cards", cards.Create)

	return mux
}
