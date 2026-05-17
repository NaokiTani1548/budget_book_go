package entity

import (
	"time"

	"github.com/google/uuid"
)

type Income struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	CategoryID   *uuid.UUID
	Amount       float64
	Description  *string
	IncomeDate   time.Time
	Memo         *string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// JOIN結果
	CategoryName *string
}