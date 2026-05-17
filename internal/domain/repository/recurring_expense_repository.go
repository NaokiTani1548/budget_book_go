package repository

import (
	"context"

	"budget-book-go/internal/domain/entity"

	"github.com/google/uuid"
)

type RecurringExpenseRepository interface {
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.RecurringExpense, error)
	FindAll(ctx context.Context, userID uuid.UUID) ([]*entity.RecurringExpense, error)
	FindActive(ctx context.Context, userID uuid.UUID) ([]*entity.RecurringExpense, error)
	Save(ctx context.Context, re *entity.RecurringExpense) (*entity.RecurringExpense, error)
	Update(ctx context.Context, re *entity.RecurringExpense) (*entity.RecurringExpense, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	ExistsLog(ctx context.Context, recurringExpenseID uuid.UUID, year int, month int) (bool, error)
	SaveLog(ctx context.Context, recurringExpenseID uuid.UUID, expenseID uuid.UUID, year int, month int) error
}