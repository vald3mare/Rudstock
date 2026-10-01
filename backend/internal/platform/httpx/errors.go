// Package httpx содержит транспортный слой HTTP: запись ответов,
// разбор запросов и перевод ошибок приложения в коды состояния.
// Это единственное место в проекте, где errs.Kind превращается в HTTP-код.
package httpx

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/vald3mare/Rudstock/backend/internal/platform/errs"
)

// ErrorResponse это тело ответа при ошибке. Наружу уходят только слаг и сообщение:
// причина (запрос к БД, ответ поставщика) остаётся в логах.
type ErrorResponse struct {
	Slug string `json:"slug"` // машинный код для фронта, например "category-not-found"
	Msg  string `json:"msg"`  // текст для пользователя
}

// statusFromKind переводит категорию ошибки в код состояния.
// default покрывает KindInternal и любую категорию, добавленную позже.
func statusFromKind(kind errs.Kind) int {
	switch kind {
	case errs.KindInvalidInput:
		return http.StatusBadRequest
	case errs.KindNotFound:
		return http.StatusNotFound
	case errs.KindForbidden:
		return http.StatusForbidden
	case errs.KindConflict:
		return http.StatusConflict
	case errs.KindTooLarge:
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusInternalServerError
	}
}

// WriteError отправляет клиенту ответ по любой ошибке и пишет её в лог.
// Хендлеры вызывают только её и не вычисляют коды состояния сами.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *errs.Error

	// errors.As, а не приведение типа: ошибку могли обернуть через %w по пути наверх.
	if errors.As(err, &appErr) {
		status := statusFromKind(appErr.Kind)

		// 5xx это наша поломка, 4xx это ошибка клиента и поводом для алерта не является.
		if status >= 500 {
			slog.Error("server error", "error", err, "method", r.Method, "path", r.URL.Path)
		} else {
			slog.Warn("client error", "error", err, "method", r.Method, "path", r.URL.Path)
		}

		WriteJSON(w, r, status, ErrorResponse{
			Slug: appErr.Slug,
			Msg:  appErr.Msg,
		})
		return
	}

	// Сюда попадают ошибки, которые не завернули в errs: из библиотек или забытые.
	// Наружу отдаём нейтральный текст, подробности только в лог.
	slog.Error("unexpected error", "error", err, "method", r.Method, "path", r.URL.Path)
	WriteJSON(w, r, http.StatusInternalServerError, ErrorResponse{
		Slug: "internal-error",
		Msg:  "Internal server error",
	})
}
