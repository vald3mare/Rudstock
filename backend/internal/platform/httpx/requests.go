package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/vald3mare/Rudstock/backend/internal/platform/errs"
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
		var maxErr *http.MaxBytesError
		if errors.As(decErr, &maxErr) {
			return errs.TooLarge("body-too-large", "Request body must not exceed 1 MB", decErr)
		}
		return errs.InvalidInput("invalid-json-body", "Invalid JSON", decErr)
	}

	// Decode читает ровно один объект. Повторный вызов должен упереться в EOF;
	// любой другой исход, включая nil, означает мусор после первого объекта.
	if !errors.Is(dec.Decode(&struct{}{}), io.EOF) {
		return errs.InvalidInput("invalid-json-body", "Request body must contain a single JSON object", nil)
	}

	return nil
}

// QueryInt64 читает целочисленный query-параметр. Если параметра нет, возвращает 0:
// что значит ноль (без фильтра, значение по умолчанию), решает вызывающий.
// Не число это ошибка клиента, а не повод молча подставить ноль.
func QueryInt64(r *http.Request, key string) (int64, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return 0, nil
	}

	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, errs.InvalidInput("invalid-query-param", "Query parameter "+key+" must be an integer", err)
	}

	return v, nil
}
