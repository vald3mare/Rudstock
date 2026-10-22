package app

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vald3mare/Rudstock/auth/internal/controllers"
	authCtrl "github.com/vald3mare/Rudstock/auth/internal/controllers/auth"
	"github.com/vald3mare/Rudstock/auth/internal/platform/config"
	sessionRepo "github.com/vald3mare/Rudstock/auth/internal/repositories/session"
	userRepo "github.com/vald3mare/Rudstock/auth/internal/repositories/user"
	authsvc "github.com/vald3mare/Rudstock/auth/internal/services/auth"
	"github.com/vald3mare/Rudstock/auth/internal/tokens"
)

// App собранное приложение auth: все слои связаны, наружу отдаётся только роутер.
type App struct {
	Router http.Handler
}

// NewApp собирает слои снизу вверх: репозитории и токены -> сервис -> хендлеры -> роутер.
// Ошибка только от NewEmitter: с неправильным секретом или сроком сервис стартовать не должен.
func NewApp(pool *pgxpool.Pool, conf *config.Config) (*App, error) {
	userRepo := userRepo.NewUserRepo(pool)
	sessionRepo := sessionRepo.NewSessionRepo(pool)

	emitter, err := tokens.NewEmitter([]byte(conf.JWTSecret), conf.JWTIssuer, conf.AccessTTL)
	if err != nil {
		return nil, fmt.Errorf("create token emitter: %w", err)
	}

	authService := authsvc.NewAuthService(userRepo, sessionRepo, emitter, conf.RefreshTTL, conf.RefreshMaxTTL)

	// один сервис уходит в оба хендлера: каждый видит его только через свой узкий интерфейс
	credentials := authCtrl.NewCredentialsHandler(authService)
	sessions := authCtrl.NewSessionHandler(authService)

	return &App{
		Router: controllers.NewRouter(credentials, sessions),
	}, nil
}
