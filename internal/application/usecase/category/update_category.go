package category

import (
	"context"

	"budget-book-go/internal/application/dto"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/domain/repository"
)

type UpdateCategoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

func NewUpdateCategoryUseCase(categoryRepo repository.CategoryRepository) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{categoryRepo: categoryRepo}
}

func (uc *UpdateCategoryUseCase) Execute(ctx context.Context, cmd dto.UpdateCategoryCommand) (*dto.CategoryResult, error) {
	if cmd.Name == "" {
		return nil, domainerror.NewInvalidInputError("カテゴリ名は必須です")
	}
	if cmd.Type != "EXPENSE" && cmd.Type != "INCOME" {
		return nil, domainerror.NewInvalidInputError("typeはEXPENSEまたはINCOMEを指定してください")
	}

	category, err := uc.categoryRepo.FindByID(ctx, cmd.ID, cmd.UserID)
	if err != nil {
		return nil, err
	}

	category.Name      = cmd.Name
	category.Type      = cmd.Type
	category.Color     = cmd.Color
	category.SortOrder = cmd.SortOrder

	updated, err := uc.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, err
	}

	return toCategoryResult(updated), nil
}