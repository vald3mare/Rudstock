package rest

import "github.com/vald3mare/Rudstock/internal/catalog/domain"

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
