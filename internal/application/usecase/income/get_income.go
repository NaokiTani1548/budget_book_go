package income

import (
	"context"
	"time"

	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/repository"

	"github.com/google/uuid"
)

type GetIncomeUseCase struct {
	incomeRepo repository.IncomeRepository
}

func NewGetIncomeUseCase(incomeRepo repository.IncomeRepository) *GetIncomeUseCase {
	return &GetIncomeUseCase{incomeRepo: incomeRepo}
}

func (uc *GetIncomeUseCase) ExecuteGetByUserID(ctx context.Context, userID uuid.UUID, limit int32, offset int32) ([]*dto.IncomeResult, int64, error) {
	totalCount, err := uc.incomeRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	incomes, err := uc.incomeRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	results := make([]*dto.IncomeResult, len(incomes))
	for i, inc := range incomes {
		results[i] = toIncomeResult(inc)
	}
	return results, totalCount, nil
}

func (uc *GetIncomeUseCase) ExecuteGetByDateRange(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]*dto.IncomeResult, error) {
	incomes, err := uc.incomeRepo.FindByDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	results := make([]*dto.IncomeResult, len(incomes))
	for i, inc := range incomes {
		results[i] = toIncomeResult(inc)
	}
	return results, nil
}

func (uc *GetIncomeUseCase) ExecuteGetOne(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.IncomeResult, error) {
	income, err := uc.incomeRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return toIncomeResult(income), nil
}

func (uc *GetIncomeUseCase) ExecuteGetPlanned(ctx context.Context, userID uuid.UUID) ([]*dto.IncomeResult, error) {
	incomes, err := uc.incomeRepo.FindPlanned(ctx, userID)
	if err != nil {
		return nil, err
	}
	results := make([]*dto.IncomeResult, len(incomes))
	for i, inc := range incomes {
		results[i] = toIncomeResult(inc)
	}
	return results, nil
}