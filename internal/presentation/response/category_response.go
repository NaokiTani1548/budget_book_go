package response

import (
	"budget-book-go/internal/application/dto"

	"github.com/google/uuid"
)

type CategoryResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Color     *string    `json:"color"`
	SortOrder int        `json:"sortOrder"`
	IsDefault bool       `json:"isDefault"`
}

func NewCategoryResponse(result *dto.CategoryResult) CategoryResponse {
	return CategoryResponse{
		ID:        result.ID,
		Name:      result.Name,
		Type:      result.Type,
		Color:     result.Color,
		SortOrder: result.SortOrder,
		IsDefault: result.IsDefault,
	}
}

func NewCategoryListResponse(results []*dto.CategoryResult) []CategoryResponse {
	responses := make([]CategoryResponse, len(results))
	for i, result := range results {
		responses[i] = NewCategoryResponse(result)
	}
	return responses
}