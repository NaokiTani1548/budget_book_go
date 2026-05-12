package repository

import (
	"budget-book-go/internal/domain/entity"
	"context"
	"time"

	"github.com/google/uuid"
)

type ExpenseRepository interface {
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Expense, error)
	FindPlanned(ctx context.Context, userID uuid.UUID) ([]*entity.Expense, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Expense, error)
	FindByDateRange(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]*entity.Expense, error)
	Save(ctx context.Context, expense *entity.Expense) (*entity.Expense, error)
	Update(ctx context.Context, expense *entity.Expense) (*entity.Expense, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
