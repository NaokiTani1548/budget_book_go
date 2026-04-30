package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCategoryCommand struct {
	UserID    uuid.UUID
	Name      string
	Type      string
	Color     *string
	SortOrder int
	IsDefault bool
}

type UpdateCategoryCommand struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Type      string
	Color     *string
	SortOrder int
}

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