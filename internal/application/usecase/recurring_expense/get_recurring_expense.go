package recurringexpense

import (
	"context"

	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/repository"

	"github.com/google/uuid"
)

type GetRecurringExpenseUseCase struct {
	recurringRepo repository.RecurringExpenseRepository
}

func NewGetRecurringExpenseUseCase(recurringRepo repository.RecurringExpenseRepository) *GetRecurringExpenseUseCase {
	return &GetRecurringExpenseUseCase{recurringRepo: recurringRepo}
}

func (uc *GetRecurringExpenseUseCase) ExecuteGetAll(ctx context.Context, userID uuid.UUID) ([]*dto.RecurringExpenseResult, error) {
	list, err := uc.recurringRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	results := make([]*dto.RecurringExpenseResult, len(list))
	for i, re := range list {
		results[i] = toRecurringExpenseResult(re)
	}
	return results, nil
}

func (uc *GetRecurringExpenseUseCase) ExecuteGetOne(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.RecurringExpenseResult, error) {
	re, err := uc.recurringRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return toRecurringExpenseResult(re), nil
}