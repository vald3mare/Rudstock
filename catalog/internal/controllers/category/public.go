package category

import (
	"context"

	"github.com/vald3mare/Rudstock/catalog/internal/models"
)

type categoryReader interface {
	List(ctx context.Context) ([]models.Category, error)
}

type PublicCategoryHandler struct {
	svc categoryReader
}

func NewPublicCategoryHandler(s categoryReader) *PublicCategoryHandler {
	return &PublicCategoryHandler{svc: s}
}
