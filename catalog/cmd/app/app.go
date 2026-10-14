package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vald3mare/Rudstock/catalog/internal/controllers"
	cardCtrl "github.com/vald3mare/Rudstock/catalog/internal/controllers/card"
	categoryCtrl "github.com/vald3mare/Rudstock/catalog/internal/controllers/category"
	cardRepo "github.com/vald3mare/Rudstock/catalog/internal/repositories/card"
	categoryRepo "github.com/vald3mare/Rudstock/catalog/internal/repositories/category"
	cardsvc "github.com/vald3mare/Rudstock/catalog/internal/services/card"
	categorySvc "github.com/vald3mare/Rudstock/catalog/internal/services/category"
)

// App собранное приложение каталога: все слои связаны, наружу отдаётся только роутер.
type App struct {
	Router http.Handler
}

// NewApp собирает слои снизу вверх: репозитории -> сервисы -> хендлеры -> роутер.
// Пул один на все репозитории: отдельный пул на каждый означал бы лишние соединения к Postgres.
func NewApp(pool *pgxpool.Pool, debugToken string) *App {
	cardRepo := cardRepo.NewCardRepo(pool)
	categoryRepo := categoryRepo.NewCategoryRepo(pool)

	cardService := cardsvc.NewCardService(cardRepo)
	categoryService := categorySvc.NewCategoryService(categoryRepo)

	cardHandler := cardCtrl.NewCardHandler(cardService)
	categoryHandler := categoryCtrl.NewCategoryHandler(categoryService)

	return &App{
		Router: controllers.NewRouter(cardHandler, categoryHandler, debugToken),
	}
}
