package dto

import (
	"time"

	"github.com/google/uuid"
)

type CategoryResult struct {
	ID        uuid.UUID
	UserID    *uuid.UUID
	Name      string
	Type      string
	Color     *string
	SortOrder int
	IsDefault bool
	CreatedAt time.Time
}