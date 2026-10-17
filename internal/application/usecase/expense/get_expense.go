package expense

import (
	"context"
	"time"

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

func (uc *GetExpenseUseCase) ExecuteGetByUserID(ctx context.Context, userID uuid.UUID, limit int32, offset int32) ([]*dto.ExpenseResult, int64, error) {
	totalCount, err := uc.expenseRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	expenses, err := uc.expenseRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	results := make([]*dto.ExpenseResult, len(expenses))
	for i, e := range expenses {
		results[i] = toExpenseResult(e)
	}
	return results, totalCount, nil
}

func (uc *GetExpenseUseCase) ExecuteGetOne(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.ExpenseResult, error) {
	expense, err := uc.expenseRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return toExpenseResult(expense), nil
}

func (uc *GetExpenseUseCase) ExecuteGetByDateRange(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]*dto.ExpenseResult, error) {
	expenses, err := uc.expenseRepo.FindByDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	results := make([]*dto.ExpenseResult, len(expenses))
	for i, e := range expenses {
		results[i] = toExpenseResult(e)
	}
	return results, nil
}

func (uc *GetExpenseUseCase) ExecuteGetPlanned(ctx context.Context, userID uuid.UUID) ([]*dto.ExpenseResult, error) {
	incomes, err := uc.expenseRepo.FindPlanned(ctx, userID)
	if err != nil {
		return nil, err
	}
	results := make([]*dto.ExpenseResult, len(incomes))
	for i, inc := range incomes {
		results[i] = toExpenseResult(inc)
	}
	return results, nil
}

func (uc *GetExpenseUseCase) ExecuteSearch(ctx context.Context, userID uuid.UUID, params repository.SearchExpenseParams) ([]*dto.ExpenseResult, error) {
	expenses, err := uc.expenseRepo.Search(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	results := make([]*dto.ExpenseResult, len(expenses))
	for i, e := range expenses {
		results[i] = toExpenseResult(e)
	}
	return results, nil
}