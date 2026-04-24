package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID
	UserID    *uuid.UUID
	Name      string
	Type      string
	Color     *string
	SortOrder int
	IsDefault bool
	CreatedAt time.Time
}