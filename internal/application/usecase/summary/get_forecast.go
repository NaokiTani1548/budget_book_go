package summary

import (
	"context"
	"time"

	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/repository"

	"github.com/google/uuid"
)

type GetForecastUseCase struct {
	summaryRepo repository.SummaryRepository
}

func NewGetForecastUseCase(summaryRepo repository.SummaryRepository) *GetForecastUseCase {
	return &GetForecastUseCase{summaryRepo: summaryRepo}
}

func (uc *GetForecastUseCase) Execute(ctx context.Context, userID uuid.UUID, targetDate time.Time) (*dto.ForecastResult, error) {
	actualIncome, err := uc.summaryRepo.SumActualIncomes(ctx, userID)
	if err != nil {
		return nil, err
	}

	actualExpense, err := uc.summaryRepo.SumActualExpenses(ctx, userID)
	if err != nil {
		return nil, err
	}

	plannedIncome, err := uc.summaryRepo.SumPlannedIncomesByDate(ctx, userID, targetDate)
	if err != nil {
		return nil, err
	}

	plannedExpense, err := uc.summaryRepo.SumPlannedExpensesByDate(ctx, userID, targetDate)
	if err != nil {
		return nil, err
	}

	currentBalance  := actualIncome - actualExpense
	forecastBalance := currentBalance + plannedIncome - plannedExpense

	return &dto.ForecastResult{
		CurrentBalance:  currentBalance,
		PlannedIncome:   plannedIncome,
		PlannedExpense:  plannedExpense,
		ForecastBalance: forecastBalance,
		TargetDate:      targetDate.Format("2006-01-02"),
	}, nil
}