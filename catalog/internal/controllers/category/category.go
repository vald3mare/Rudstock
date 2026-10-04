package category

import (
	"github.com/vald3mare/Rudstock/catalog/internal/services/category"
)

// CategoryHandler HTTP-ручки категорий: разбирает запрос, зовёт сервис, пишет ответ.
type CategoryHandler struct {
	svc category.CategoryService
}

// NewCategoryHandler конструктор хендлера категорий.
func NewCategoryHandler(s category.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: s}
}
