package category

import (
	"context"

	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/entity"
	"budget-book-go/internal/domain/repository"

	"github.com/google/uuid"
)

type GetCategoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

func NewGetCategoryUseCase(categoryRepo repository.CategoryRepository) *GetCategoryUseCase {
	return &GetCategoryUseCase{categoryRepo: categoryRepo}
}

func (uc *GetCategoryUseCase) ExecuteGetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.CategoryResult, error) {
	categories, err := uc.categoryRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	results := make([]*dto.CategoryResult, len(categories))
	for i, c := range categories {
		results[i] = toCategoryResult(c)
	}
	return results, nil
}

func toCategoryResult(c *entity.Category) *dto.CategoryResult {
	return &dto.CategoryResult{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		Type:      c.Type,
		Color:     c.Color,
		SortOrder: c.SortOrder,
		IsDefault: c.IsDefault,
		CreatedAt: c.CreatedAt,
	}
}