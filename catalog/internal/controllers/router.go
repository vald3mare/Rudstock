package controllers

import (
	"net/http"

	"github.com/vald3mare/Rudstock/catalog/internal/controllers/card"
	"github.com/vald3mare/Rudstock/catalog/internal/controllers/category"
	"github.com/vald3mare/Rudstock/catalog/internal/infrastructure/middleware/auth"
)

func NewRouter(cards *card.CardHandler, categories *category.CategoryHandler) http.Handler {
	mux := http.NewServeMux()
	adminMux := http.NewServeMux()

	// public
	mux.HandleFunc("GET /categories", categories.List)
	mux.HandleFunc("GET /cards", cards.List)
	mux.HandleFunc("GET /cards/{id}", cards.Get)

	// admin
	adminMux.HandleFunc("POST /admin/cards", cards.Create)
	adminMux.HandleFunc("PATCH /admin/cards/{id}", cards.Update)
	adminMux.HandleFunc("DELETE /admin/cards/{id}", cards.Delete)
	adminMux.HandleFunc("POST /admin/categories", categories.Create)
	adminMux.HandleFunc("PATCH /admin/categories/{id}", categories.Update)
	adminMux.HandleFunc("DELETE /admin/categories/{id}", categories.Delete)

	// всё под /admin/ сначала проходит авторизацию, потом ищет маршрут в adminMux;
	// пути в adminMux полные, поэтому StripPrefix не нужен. Публичные ручки без неё
	mux.Handle("/admin/", auth.Middleware(adminMux))

	return mux
}
