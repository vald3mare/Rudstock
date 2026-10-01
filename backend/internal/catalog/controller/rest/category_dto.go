package rest

import (
	"github.com/vald3mare/Rudstock/backend/internal/catalog/domain"
	"github.com/vald3mare/Rudstock/backend/internal/catalog/service"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

// UpdateCategoryRequest тело PATCH: указатель, чтобы отличить отсутствующее поле от пустой строки.
type UpdateCategoryRequest struct {
	Name *string `json:"name"`
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
		resp = append(resp, toCategoryResponse(c))
	}
	return resp
}

// toCategoryResponse маппинг одной категории, единственное место маппинга полей.
func toCategoryResponse(c domain.Category) CategoryResponse {
	return CategoryResponse{ID: c.ID, Name: c.Name}
}

// toInput переводит запрос по контракту во входные данные сервиса
func (r CreateCategoryRequest) toInput() service.CreateCategoryInput {
	return service.CreateCategoryInput{
		Name: r.Name,
	}
}

// toInput переводит тело PATCH во входные данные сервиса, указатели передаются как есть
func (r UpdateCategoryRequest) toInput() service.UpdateCategoryInput {
	return service.UpdateCategoryInput{
		Name: r.Name,
	}
}
