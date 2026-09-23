// Package errs описывает ошибки приложения в терминах предметной области,
// а не транспорта. Слои app и domain возвращают ошибки этого пакета;
// перевод в HTTP-коды делает только пакет httpx.
package errs

// Kind это категория ошибки. Категории намеренно не привязаны к HTTP:
// те же ошибки возвращаются из воркеров и CLI, где кодов статуса нет.
type Kind string

const (
	KindInvalidInput Kind = "invalid_input"
	KindNotFound     Kind = "not_found"
	KindForbidden    Kind = "forbidden"
	KindConflict     Kind = "conflict"
	KindInternal     Kind = "internal"
)

// Error переносит ошибку между слоями.
// Наружу отдаются только Slug и Msg; err содержит исходную причину
// (запрос к БД, ответ поставщика) и предназначена исключительно для логов.
type Error struct {
	Kind Kind   // категория, по ней httpx выбирает код ответа
	Slug string // стабильный машинный код, например "product-not-found"
	Msg  string // текст для пользователя
	err  error  // причина, неэкспортируемая, чтобы не утекла в ответ
}

// Error возвращает максимум информации для лога: слаг и причину, если она есть.
func (e *Error) Error() string {
	if e.err != nil {
		return e.Slug + ": " + e.err.Error()
	}
	return e.Slug + ": " + e.Msg
}

// Unwrap открывает причину для errors.Is и errors.As.
func (e *Error) Unwrap() error {
	return e.err
}

// Конструкторы ниже проставляют категорию сами, поэтому перепутать её нельзя.
// Аргумент err может быть nil, когда ошибку породило само правило,
// а не вызов внешней системы.

func NotFound(slug string, msg string, err error) *Error {
	return &Error{
		Kind: KindNotFound,
		Slug: slug,
		Msg:  msg,
		err:  err,
	}
}

func InvalidInput(slug string, msg string, err error) *Error {
	return &Error{
		Kind: KindInvalidInput,
		Slug: slug,
		Msg:  msg,
		err:  err,
	}
}

func Forbidden(slug string, msg string, err error) *Error {
	return &Error{
		Kind: KindForbidden,
		Slug: slug,
		Msg:  msg,
		err:  err,
	}
}

func Conflict(slug string, msg string, err error) *Error {
	return &Error{
		Kind: KindConflict,
		Slug: slug,
		Msg:  msg,
		err:  err,
	}
}

func Internal(slug string, msg string, err error) *Error {
	return &Error{
		Kind: KindInternal,
		Slug: slug,
		Msg:  msg,
		err:  err,
	}
}
