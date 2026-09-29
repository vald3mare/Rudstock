package rest

import "github.com/google/uuid"

type CreateCardRequest struct {
	CategoryID  int64  `json:"category_id"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	PhotoURL    string `json:"photo_url"`
}

type CreateCardResponse struct {
	CardID uuid.UUID `json:"card_id"`
}
