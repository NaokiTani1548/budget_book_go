package entity

import (
	"time"
	"github.com/google/uuid"
)

type Expense struct {
	ID uuid.UUID
	UserID uuid.UUID
	CategoryID *uuid.UUID
	Amount float64
	Description *string
	ExpenseDate time.Time
	PaymentMethod *string
	Memo *string
	CreatedAt time.Time
	UpdatedAt time.Time

	CategoryName *string
}