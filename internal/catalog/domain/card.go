package domain

import (
	"time"

	"github.com/google/uuid"
)

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
