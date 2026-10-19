package controllers

import (
	"net/http"

	"github.com/vald3mare/Rudstock/catalog/internal/controllers/card"
	"github.com/vald3mare/Rudstock/catalog/internal/controllers/category"
	"github.com/vald3mare/Rudstock/catalog/internal/infrastructure/middleware/auth"
)

func NewRouter(
	publicCards *card.PublicCardHandler,
	adminCards *card.AdminCardHandler,
	publicCategories *category.PublicCategoryHandler,
	adminCategories *category.AdminCategoryHandler,
	debugToken string,
) http.Handler {
	mux := http.NewServeMux()
	adminMux := http.NewServeMux()

	// public
	mux.HandleFunc("GET /categories", publicCategories.List)
	mux.HandleFunc("GET /cards", publicCards.List)
	mux.HandleFunc("GET /cards/{id}", publicCards.Get)

	// admin
	adminMux.HandleFunc("POST /admin/cards", adminCards.Create)
	adminMux.HandleFunc("PATCH /admin/cards/{id}", adminCards.Update)
	adminMux.HandleFunc("DELETE /admin/cards/{id}", adminCards.Delete)
	adminMux.HandleFunc("POST /admin/categories", adminCategories.Create)
	adminMux.HandleFunc("PATCH /admin/categories/{id}", adminCategories.Update)
	adminMux.HandleFunc("DELETE /admin/categories/{id}", adminCategories.Delete)

	// всё под /admin/ сначала проходит авторизацию, потом ищет маршрут в adminMux;
	// пути в adminMux полные, поэтому StripPrefix не нужен. Публичные ручки без неё
	mux.Handle("/admin/", auth.Middleware(debugToken, adminMux))

	return mux
}
