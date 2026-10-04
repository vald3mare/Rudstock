package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vald3mare/Rudstock/catalog/internal/controllers"
	cardctrl "github.com/vald3mare/Rudstock/catalog/internal/controllers/card"
	categoryctrl "github.com/vald3mare/Rudstock/catalog/internal/controllers/category"
	cardrepo "github.com/vald3mare/Rudstock/catalog/internal/repositories/card"
	categoryrepo "github.com/vald3mare/Rudstock/catalog/internal/repositories/category"
	cardsvc "github.com/vald3mare/Rudstock/catalog/internal/services/card"
	categorysvc "github.com/vald3mare/Rudstock/catalog/internal/services/category"
)

// App собранное приложение каталога: все слои связаны, наружу отдаётся только роутер.
type App struct {
	Router http.Handler
}

// NewApp собирает слои снизу вверх: репозитории -> сервисы -> хендлеры -> роутер.
// Пул один на все репозитории: отдельный пул на каждый означал бы лишние соединения к Postgres.
func NewApp(pool *pgxpool.Pool) *App {
	cardRepo := cardrepo.NewCardRepo(pool)
	categoryRepo := categoryrepo.NewCategoryRepo(pool)

	cardService := cardsvc.NewCardService(cardRepo)
	categoryService := categorysvc.NewCategoryService(categoryRepo)

	cardHandler := cardctrl.NewCardHandler(cardService)
	categoryHandler := categoryctrl.NewCategoryHandler(categoryService)

	return &App{
		Router: controllers.NewRouter(cardHandler, categoryHandler),
	}
}
