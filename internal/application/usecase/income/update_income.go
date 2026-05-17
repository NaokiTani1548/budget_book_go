package income

import (
	"context"

	"budget-book-go/internal/application/dto"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/domain/repository"
	"budget-book-go/internal/domain/valueobject"
)

type UpdateIncomeUseCase struct {
	incomeRepo repository.IncomeRepository
}

func NewUpdateIncomeUseCase(incomeRepo repository.IncomeRepository) *UpdateIncomeUseCase {
	return &UpdateIncomeUseCase{incomeRepo: incomeRepo}
}

func (uc *UpdateIncomeUseCase) Execute(ctx context.Context, cmd dto.UpdateIncomeCommand) (*dto.IncomeResult, error) {
	_, err := valueobject.NewMoney(cmd.Amount)
	if err != nil {
		return nil, domainerror.NewInvalidInputError(err.Error())
	}

	income, err := uc.incomeRepo.FindByID(ctx, cmd.ID, cmd.UserID)
	if err != nil {
		return nil, err
	}

	income.CategoryID  = cmd.CategoryID
	income.Amount      = cmd.Amount
	income.Description = cmd.Description
	income.IncomeDate  = cmd.IncomeDate
	income.Memo        = cmd.Memo

	updated, err := uc.incomeRepo.Update(ctx, income)
	if err != nil {
		return nil, err
	}
	return toIncomeResult(updated), nil
}