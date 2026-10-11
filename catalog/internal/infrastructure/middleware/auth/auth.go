package auth

import (
	"net/http"

	"github.com/vald3mare/Rudstock/catalog/internal/platform/httpx"
)

// Заглушка до появления сервиса auth: пропускает запрос только с тестовым заголовком.
// Потом здесь будет проверка Bearer-токена и роли сотрудника.
const (
	debugTokenHeader = "X-Debug-Token"
	debugToken       = "dev"
)

// Middleware пускает дальше только запросы с X-Debug-Token: dev, остальным отвечает 401.
// Вешается на /admin/*, публичные ручки её не проходят.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(debugTokenHeader) != debugToken {
			// отказ: пишем ответ и выходим, next не вызываем,
			// иначе хендлер допишет свой JSON в тот же ответ
			httpx.WriteJSON(w, r, http.StatusUnauthorized, httpx.ErrorResponse{
				Slug: "unauthorized",
				Msg:  "Authorization required",
			})
			return
		}

		// пропуск: сами ничего не пишем, ответ формирует хендлер
		next.ServeHTTP(w, r)
	})
}
