package controllers

import (
	"net/http"

	"github.com/vald3mare/Rudstock/catalog/internal/controllers/card"
	"github.com/vald3mare/Rudstock/catalog/internal/controllers/category"
)

func NewRouter(cards *card.CardHandler, categories *category.CategoryHandler) http.Handler {
	mux := http.NewServeMux()

	// public
	mux.HandleFunc("GET /categories", categories.List)
	mux.HandleFunc("GET /cards", cards.List)
	mux.HandleFunc("GET /cards/{id}", cards.Get)

	// admin
	mux.HandleFunc("POST /admin/cards", cards.Create)
	mux.HandleFunc("PATCH /admin/cards/{id}", cards.Update)
	mux.HandleFunc("DELETE /admin/cards/{id}", cards.Delete)
	mux.HandleFunc("POST /admin/categories", categories.Create)
	mux.HandleFunc("PATCH /admin/categories/{id}", categories.Update)
	mux.HandleFunc("DELETE /admin/categories/{id}", categories.Delete)

	return mux
}
