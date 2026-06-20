package repository

import (
	"budget-book-go/internal/domain/entity"
	"context"
	"time"

	"github.com/google/uuid"
)
type SearchExpenseParams struct {
	DateFrom   *time.Time
	DateTo     *time.Time
	CategoryID *uuid.UUID
	Keyword    *string
}

type ExpenseRepository interface {
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Expense, error)
	FindPlanned(ctx context.Context, userID uuid.UUID) ([]*entity.Expense, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Expense, error)
	FindByDateRange(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]*entity.Expense, error)
	Search(ctx context.Context, userID uuid.UUID, params SearchExpenseParams) ([]*entity.Expense, error)
	Save(ctx context.Context, expense *entity.Expense) (*entity.Expense, error)
	Update(ctx context.Context, expense *entity.Expense) (*entity.Expense, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
