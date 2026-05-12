package expense

import (
	"context"
	"log"

	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/entity"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/domain/repository"
	"budget-book-go/internal/domain/valueobject"
)

type CreateExpenseUseCase struct {
	expenseRepo repository.ExpenseRepository
}

func NewCreateExpenseUseCase(expenseRepo repository.ExpenseRepository) *CreateExpenseUseCase {
	return &CreateExpenseUseCase{expenseRepo: expenseRepo}
}

func (uc *CreateExpenseUseCase) Execute(ctx context.Context, cmd dto.CreateExpenseCommand) (*dto.ExpenseResult, error) {
	_, err := valueobject.NewMoney(cmd.Amount)
	if err != nil {
		return nil, domainerror.NewInvalidInputError(err.Error())
	}

	log.Printf("cmd.IsPlanned: %v, cmd.PlannedDate: %v", cmd.IsPlanned, cmd.PlannedDate)
	expense := &entity.Expense{
		UserID:        cmd.UserID,
		CategoryID:    cmd.CategoryID,
		Amount:        cmd.Amount,
		Description:   cmd.Description,
		ExpenseDate:   cmd.ExpenseDate,
		PaymentMethod: cmd.PaymentMethod,
		Memo:          cmd.Memo,
		IsPlanned:     cmd.IsPlanned,
		PlannedDate:   cmd.PlannedDate,
	}
	saved, err := uc.expenseRepo.Save(ctx, expense)
	if err != nil {
		return nil, err
	}

	return toExpenseResult(saved), nil
}