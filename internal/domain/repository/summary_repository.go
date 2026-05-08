package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SummaryRepository interface {
	SumActualIncomes(ctx context.Context, userID uuid.UUID) (float64, error)
	SumActualExpenses(ctx context.Context, userID uuid.UUID) (float64, error)
	SumPlannedIncomesByDate(ctx context.Context, userID uuid.UUID, targetDate time.Time) (float64, error)
	SumPlannedExpensesByDate(ctx context.Context, userID uuid.UUID, targetDate time.Time) (float64, error)
}