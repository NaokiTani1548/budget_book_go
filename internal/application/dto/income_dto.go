package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateIncomeCommand struct {
	UserID      uuid.UUID
	CategoryID  *uuid.UUID
	Amount      float64
	Description *string
	IncomeDate  time.Time
	Memo        *string
}

type UpdateIncomeCommand struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	CategoryID  *uuid.UUID
	Amount      float64
	Description *string
	IncomeDate  time.Time
	Memo        *string
}

type IncomeResult struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	CategoryID  *uuid.UUID
	CategoryName *string
	Amount      float64
	Description *string
	IncomeDate  time.Time
	Memo        *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}