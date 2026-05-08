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
	IsPlanned   bool
	PlannedDate *time.Time
}

type UpdateIncomeCommand struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	CategoryID  *uuid.UUID
	Amount      float64
	Description *string
	IncomeDate  time.Time
	Memo        *string
	IsPlanned   bool
	PlannedDate *time.Time
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
	IsPlanned   bool
	PlannedDate *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}