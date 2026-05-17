package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateRecurringExpenseCommand struct {
	UserID        uuid.UUID
	CategoryID    *uuid.UUID
	Amount        float64
	Description   *string
	PaymentMethod *string
	Memo          *string
	BillingDay    int
	StartDate     time.Time
	EndDate       *time.Time
}

type UpdateRecurringExpenseCommand struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	CategoryID    *uuid.UUID
	Amount        float64
	Description   *string
	PaymentMethod *string
	Memo          *string
	BillingDay    int
	StartDate     time.Time
	EndDate       *time.Time
	IsActive      bool
}

type RecurringExpenseResult struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	CategoryID    *uuid.UUID
	CategoryName  *string
	Amount        float64
	Description   *string
	PaymentMethod *string
	Memo          *string
	BillingDay    int
	StartDate     time.Time
	EndDate       *time.Time
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}