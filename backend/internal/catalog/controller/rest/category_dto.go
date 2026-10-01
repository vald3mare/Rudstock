package rest

import (
	"github.com/vald3mare/Rudstock/backend/internal/catalog/domain"
	"github.com/vald3mare/Rudstock/backend/internal/catalog/service"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type CreateCategoryResponse struct {
	CategoryID int64 `json:"category_id"`
}

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// toCategoryResponses отделяет JSON-формат от доменной модели:
// переименование поля в domain не должно менять ответ API.
func toCategoryResponses(categories []domain.Category) []CategoryResponse {
	// make, а не var: пустой список должен уйти как [], а не null
	resp := make([]CategoryResponse, 0, len(categories))
	for _, c := range categories {
		resp = append(resp, CategoryResponse{ID: c.ID, Name: c.Name})
	}
	return resp
}

// toInput переводит запрос по контракту во входные данные сервиса
func (r CreateCategoryRequest) toInput() service.CreateCategoryInput {
	return service.CreateCategoryInput{
		Name: r.Name,
	}
}
