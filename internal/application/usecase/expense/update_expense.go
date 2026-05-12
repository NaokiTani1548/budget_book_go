package expense

import (
	"context"

	"budget-book-go/internal/application/dto"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/domain/repository"
	"budget-book-go/internal/domain/valueobject"
)

type UpdateExpenseUseCase struct {
	expenseRepo repository.ExpenseRepository
}

func NewUpdateExpenseUseCase(expenseRepo repository.ExpenseRepository) *UpdateExpenseUseCase {
	return &UpdateExpenseUseCase{expenseRepo: expenseRepo}
}

func (uc *UpdateExpenseUseCase) Execute(ctx context.Context, cmd dto.UpdateExpenseCommand) (*dto.ExpenseResult, error) {
	// 金額バリデーション
	_, err := valueobject.NewMoney(cmd.Amount)
	if err != nil {
		return nil, domainerror.NewInvalidInputError(err.Error())
	}

	// 対象データの存在確認
	expense, err := uc.expenseRepo.FindByID(ctx, cmd.ID, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// 更新
	expense.CategoryID    = cmd.CategoryID
	expense.Amount        = cmd.Amount
	expense.Description   = cmd.Description
	expense.ExpenseDate   = cmd.ExpenseDate
	expense.PaymentMethod = cmd.PaymentMethod
	expense.Memo          = cmd.Memo
	expense.IsPlanned     = cmd.IsPlanned
	expense.PlannedDate   = cmd.PlannedDate

	updated, err := uc.expenseRepo.Update(ctx, expense)
	if err != nil {
		return nil, err
	}

	return toExpenseResult(updated), nil
}