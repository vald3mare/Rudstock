package card

import (
	"context"

	"github.com/google/uuid"
	"github.com/vald3mare/Rudstock/catalog/internal/models"
	cardsvc "github.com/vald3mare/Rudstock/catalog/internal/services/card"
)

type cardWriter interface {
	Create(ctx context.Context, in cardsvc.CreateCardInput) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, in cardsvc.UpdateCardInput) (models.Card, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type AdminCardHandler struct {
	svc cardWriter
}

func NewAdminCardHandler(s cardWriter) *AdminCardHandler {
	return &AdminCardHandler{svc: s}
}
