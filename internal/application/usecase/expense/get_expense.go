package expense

import (
	"context"

	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/repository"

	"github.com/google/uuid"
)

type GetExpenseUseCase struct {
	expenseRepo repository.ExpenseRepository
}

func NewGetExpenseUseCase(expenseRepo repository.ExpenseRepository) *GetExpenseUseCase {
	return &GetExpenseUseCase{expenseRepo: expenseRepo}
}

func (uc *GetExpenseUseCase) ExecuteGetByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.ExpenseResult, error) {
	expenses, err := uc.expenseRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	results := make([]*dto.ExpenseResult, len(expenses))
	for i, e := range expenses {
		results[i] = toExpenseResult(e)
	}
	return results, nil
}

func (uc *GetExpenseUseCase) ExecuteGetOne(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.ExpenseResult, error) {
	expense, err := uc.expenseRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return toExpenseResult(expense), nil
}