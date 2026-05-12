package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateExpenseCommand struct {
	UserID        uuid.UUID
	CategoryID    *uuid.UUID
	Amount        float64
	Description   *string
	ExpenseDate   time.Time
	PaymentMethod *string
	Memo          *string
	IsPlanned   bool
	PlannedDate *time.Time
}

type UpdateExpenseCommand struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	CategoryID    *uuid.UUID
	Amount        float64
	Description   *string
	ExpenseDate   time.Time
	PaymentMethod *string
	Memo          *string
	IsPlanned   bool
	PlannedDate *time.Time
}

type ExpenseResult struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	CategoryID    *uuid.UUID
	CategoryName  *string
	Amount        float64
	Description   *string
	ExpenseDate   time.Time
	PaymentMethod *string
	Memo          *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	IsPlanned   bool
	PlannedDate *time.Time
}

