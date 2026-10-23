package category

import (
	"context"

	"github.com/vald3mare/Rudstock/catalog/internal/models"
	categorysvc "github.com/vald3mare/Rudstock/catalog/internal/services/category"
)

type categoryWriter interface {
	Create(ctx context.Context, in categorysvc.CreateCategoryInput) (int64, error)
	Update(ctx context.Context, id int64, in categorysvc.UpdateCategoryInput) (models.Category, error)
	Delete(ctx context.Context, id int64) error
}

type AdminCategoryHandler struct {
	svc categoryWriter
}

func NewAdminCategoryHandler(s categoryWriter) *AdminCategoryHandler {
	return &AdminCategoryHandler{svc: s}
}
