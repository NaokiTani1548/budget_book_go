package postgres

import (
	"context"
	"time"

	"budget-book-go/internal/domain/repository"
	dbsqlc "budget-book-go/internal/infrastructure/persistence/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type summaryRepository struct {
	db      *pgxpool.Pool
	queries *dbsqlc.Queries
}

func NewSummaryRepository(db *pgxpool.Pool) repository.SummaryRepository {
	return &summaryRepository{
		db:      db,
		queries: dbsqlc.New(db),
	}
}

func (r *summaryRepository) SumActualIncomes(ctx context.Context, userID uuid.UUID) (float64, error) {
	total, err := r.queries.SumActualIncomes(ctx, uuidToPgtype(userID))
	if err != nil {
		return 0, err
	}
	return numericToFloat(total), nil
}

func (r *summaryRepository) SumActualExpenses(ctx context.Context, userID uuid.UUID) (float64, error) {
	total, err := r.queries.SumActualExpenses(ctx, uuidToPgtype(userID))
	if err != nil {
		return 0, err
	}
	return numericToFloat(total), nil
}

func (r *summaryRepository) SumPlannedIncomesByDate(ctx context.Context, userID uuid.UUID, targetDate time.Time) (float64, error) {
	total, err := r.queries.SumPlannedIncomesByDate(ctx, dbsqlc.SumPlannedIncomesByDateParams{
		UserID:     uuidToPgtype(userID),
		IncomeDate: dateToPgtype(targetDate),
	})
	if err != nil {
		return 0, err
	}
	return numericToFloat(total), nil
}

func (r *summaryRepository) SumPlannedExpensesByDate(ctx context.Context, userID uuid.UUID, targetDate time.Time) (float64, error) {
	total, err := r.queries.SumPlannedExpensesByDate(ctx, dbsqlc.SumPlannedExpensesByDateParams{
		UserID:      uuidToPgtype(userID),
		ExpenseDate: dateToPgtype(targetDate),
	})
	if err != nil {
		return 0, err
	}
	return numericToFloat(total), nil
}