package rest

import "net/http"

func NewRouter(cards *CardHandler, categories *CategoryHandler) http.Handler {
	mux := http.NewServeMux()

	// public
	mux.HandleFunc("GET /categories", categories.List)
	mux.HandleFunc("GET /cards", cards.List)
	mux.HandleFunc("GET /cards/{id}", cards.Get)

	// admin
	mux.HandleFunc("POST /admin/cards", cards.Create)
	mux.HandleFunc("POST /admin/categories", categories.Create)

	return mux
}
