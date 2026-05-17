package income

import (
	"context"

	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/entity"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/domain/repository"
	"budget-book-go/internal/domain/valueobject"
)

type CreateIncomeUseCase struct {
	incomeRepo repository.IncomeRepository
}

func NewCreateIncomeUseCase(incomeRepo repository.IncomeRepository) *CreateIncomeUseCase {
	return &CreateIncomeUseCase{incomeRepo: incomeRepo}
}

func (uc *CreateIncomeUseCase) Execute(ctx context.Context, cmd dto.CreateIncomeCommand) (*dto.IncomeResult, error) {
	_, err := valueobject.NewMoney(cmd.Amount)
	if err != nil {
		return nil, domainerror.NewInvalidInputError(err.Error())
	}

	income := &entity.Income{
		UserID:      cmd.UserID,
		CategoryID:  cmd.CategoryID,
		Amount:      cmd.Amount,
		Description: cmd.Description,
		IncomeDate:  cmd.IncomeDate,
		Memo:        cmd.Memo,
	}

	saved, err := uc.incomeRepo.Save(ctx, income)
	if err != nil {
		return nil, err
	}
	return toIncomeResult(saved), nil
}