package category

import (
	"context"

	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/entity"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/domain/repository"
)

type CreateCategoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

func NewCreateCategoryUseCase(categoryRepo repository.CategoryRepository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{categoryRepo: categoryRepo}
}

func (uc *CreateCategoryUseCase) Execute(ctx context.Context, cmd dto.CreateCategoryCommand) (*dto.CategoryResult, error) {
	if cmd.Name == "" {
		return nil, domainerror.NewInvalidInputError("カテゴリ名は必須です")
	}
	if cmd.Type != "EXPENSE" && cmd.Type != "INCOME" {
		return nil, domainerror.NewInvalidInputError("typeはEXPENSEまたはINCOMEを指定してください")
	}

	category := &entity.Category{
		UserID:    &cmd.UserID,
		Name:      cmd.Name,
		Type:      cmd.Type,
		Color:     cmd.Color,
		SortOrder: cmd.SortOrder,
		IsDefault: cmd.IsDefault,
	}

	saved, err := uc.categoryRepo.Save(ctx, category)
	if err != nil {
		return nil, err
	}

	return toCategoryResult(saved), nil
}