package cors

import "net/http"

// Браузер шлёт заголовок Origin на запросы с другого адреса (фронт :5173 -> catalog :8081)
// и не отдаёт ответ коду страницы, если сервер явно не разрешил этот Origin.
const (
	allowedMethods = "GET, POST, PATCH, DELETE"
	allowedHeaders = "Content-Type, Authorization, X-Debug-Token"
	maxAge         = "600" // сколько секунд браузер кэширует ответ на preflight
)

/*CORS middleware добавляет заголовки: Access-Control-Allow-Origin, Allow-Methods, Allow-Headers, Allow-Credentials.
В Go — rs/cors или свой middleware. Для preflight (OPTIONS) — отвечает 204 с нужными headers.
Настройки: AllowedOrigins (список доменов, НЕ * если credentials), AllowedMethods, AllowedHeaders, MaxAge (кеш preflight).
В chi: r.Use(cors.Handler(cors.Options{...})).
Не ставь Allow-Origin: * с credentials — браузер заблокирует.
*/

// func Middleware(allowed []string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {

// 	}
// }

func Middleware(allowed []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}
