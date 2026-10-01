package rest

import (
	"time"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/internal/catalog/domain"
	"github.com/vald3mare/Rudstock/internal/catalog/service"
)

type CreateCardRequest struct {
	CategoryID  int64  `json:"category_id"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	PhotoURL    string `json:"photo_url"`
}

type CreateCardResponse struct {
	CardID uuid.UUID `json:"card_id"`
}

type CardResponse struct {
	ID          uuid.UUID `json:"id"`
	CategoryID  int64     `json:"category_id"`
	Description string    `json:"description"`
	Price       int64     `json:"price"` // копейки
	PhotoURL    string    `json:"photo_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListCardsResponse struct {
	Items []CardResponse `json:"items"`
	Total int64          `json:"total"`
}

// toCardResponses отделяет JSON-формат от доменной модели.
func toCardResponses(cards []domain.Card) []CardResponse {
	// make, а не var: пустой список должен уйти как [], а не null
	resp := make([]CardResponse, 0, len(cards))
	for _, c := range cards {
		resp = append(resp, CardResponse{
			ID:          c.ID,
			CategoryID:  c.CategoryID,
			Description: c.Description,
			Price:       c.Price,
			PhotoURL:    c.PhotoURL,
			CreatedAt:   c.CreatedAt,
		})
	}
	return resp
}

// toInput переводит запрос по контракту во входные данные сервиса
// метод только читает поля, а не модфицирует, поэтому по значению передаем ресивера
func (r CreateCardRequest) toInput() service.CreateCardInput {
	return service.CreateCardInput{
		CategoryID:  r.CategoryID,
		Description: r.Description,
		Price:       r.Price,
		PhotoURL:    r.PhotoURL,
	}
}
