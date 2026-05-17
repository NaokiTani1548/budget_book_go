package recurringexpense

import (
	"context"

	"budget-book-go/internal/domain/repository"

	"github.com/google/uuid"
)

type DeleteRecurringExpenseUseCase struct {
	recurringRepo repository.RecurringExpenseRepository
}

func NewDeleteRecurringExpenseUseCase(recurringRepo repository.RecurringExpenseRepository) *DeleteRecurringExpenseUseCase {
	return &DeleteRecurringExpenseUseCase{recurringRepo: recurringRepo}
}

func (uc *DeleteRecurringExpenseUseCase) Execute(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := uc.recurringRepo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return uc.recurringRepo.Delete(ctx, id, userID)
}