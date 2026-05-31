package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash *string
	Provider     string
	ProviderID   *string
	Name         *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}