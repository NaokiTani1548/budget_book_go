package expense

import (
	"context"

	"budget-book-go/internal/domain/repository"

	"github.com/google/uuid"
)

type DeleteExpenseUseCase struct {
	expenseRepo repository.ExpenseRepository
}

func NewDeleteExpenseUseCase(expenseRepo repository.ExpenseRepository) *DeleteExpenseUseCase {
	return &DeleteExpenseUseCase{expenseRepo: expenseRepo}
}

func (uc *DeleteExpenseUseCase) Execute(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// 対象データの存在確認
	_, err := uc.expenseRepo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}

	return uc.expenseRepo.Delete(ctx, id, userID)
}