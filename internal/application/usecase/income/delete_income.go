package income

import (
	"context"

	"budget-book-go/internal/domain/repository"

	"github.com/google/uuid"
)

type DeleteIncomeUseCase struct {
	incomeRepo repository.IncomeRepository
}

func NewDeleteIncomeUseCase(incomeRepo repository.IncomeRepository) *DeleteIncomeUseCase {
	return &DeleteIncomeUseCase{incomeRepo: incomeRepo}
}

func (uc *DeleteIncomeUseCase) Execute(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := uc.incomeRepo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return uc.incomeRepo.Delete(ctx, id, userID)
}