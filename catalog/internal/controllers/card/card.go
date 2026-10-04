package card

import (
	"github.com/vald3mare/Rudstock/catalog/internal/services/card"
)

// CardHandler HTTP-ручки карточек: разбирает запрос, зовёт сервис, пишет ответ.
type CardHandler struct {
	svc card.CardService
}

// NewCardHandler конструктор хендлера карточек.
func NewCardHandler(s card.CardService) *CardHandler {
	return &CardHandler{svc: s}
}
