package entity

import (
	"time"

	"github.com/google/uuid"
)

type RecurringExpense struct {
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
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// JOIN結果
	CategoryName *string
}