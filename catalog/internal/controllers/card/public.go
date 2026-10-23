package card

import (
	"context"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
	cardsvc "github.com/vald3mare/Rudstock/catalog/internal/services/card"
)

type cardReader interface {
	List(ctx context.Context, in cardsvc.ListCardsInput) ([]models.Card, int64, error)
	Get(ctx context.Context, id uuid.UUID) (models.Card, error)
}

type PublicCardHandler struct {
	svc cardReader
}

func NewPublicCardHandler(s cardReader) *PublicCardHandler {
	return &PublicCardHandler{svc: s}
}
