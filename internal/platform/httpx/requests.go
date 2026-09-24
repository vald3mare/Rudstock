package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/vald3mare/Rudstock/internal/platform/errs"
)

// DecodeJSON разбирает тело запроса в dst и возвращает ошибку приложения,
// если тело некорректно. Ответ клиенту не пишет: это делает хендлер
// через WriteError. Текст ошибки от encoding/json наружу не уходит,
// он попадает только в лог.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Ограничение защищает от тела произвольного размера, которое иначе
	// прочиталось бы в память целиком. w нужен MaxBytesReader, чтобы
	// корректно оборвать соединение при превышении.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB

	dec := json.NewDecoder(r.Body)

	// Без этого encoding/json молча игнорирует лишние поля:
	// опечатка в имени поля прошла бы как валидный запрос с нулевым значением.
	dec.DisallowUnknownFields()

	decErr := dec.Decode(dst)
	if decErr != nil {
		return errs.InvalidInput("invalid-json-body", "Invalid JSON", decErr)
	}

	// Decode читает ровно один объект. Повторный вызов должен упереться в EOF;
	// любой другой исход, включая nil, означает мусор после первого объекта.
	if !errors.Is(dec.Decode(&struct{}{}), io.EOF) {
		return errs.InvalidInput("invalid-json-body", "Request body must contain a single JSON object", nil)
	}

	return nil
}
