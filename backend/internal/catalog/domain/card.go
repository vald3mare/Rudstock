package domain

import (
	"time"

	"github.com/google/uuid"
)

// Card карточка товара.
type Card struct {
	ID          uuid.UUID
	CategoryID  int64
	Description string
	Price       int64 // в копейках
	PhotoURL    string
	CreatedAt   time.Time
}

// CardFilter условия выборки списка карточек.
type CardFilter struct {
	CategoryID int64 // 0 означает без фильтра по категории
	Limit      int
	Offset     int
}

// CardPatch частичное обновление карточки.
// Указатели нужны, чтобы отличить "поле не передали" (nil, не меняем)
// от "передали нулевое значение" (например, пустое описание).
type CardPatch struct {
	CategoryID  *int64
	Description *string
	Price       *int64 // в копейках
	PhotoURL    *string
}
