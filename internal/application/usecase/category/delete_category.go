package category

import (
	"context"

	"budget-book-go/internal/domain/repository"

	"github.com/google/uuid"
)

type DeleteCategoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

func NewDeleteCategoryUseCase(categoryRepo repository.CategoryRepository) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{categoryRepo: categoryRepo}
}

func (uc *DeleteCategoryUseCase) Execute(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := uc.categoryRepo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return uc.categoryRepo.Delete(ctx, id, userID)
}