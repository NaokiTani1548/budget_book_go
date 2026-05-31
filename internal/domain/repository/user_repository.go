package repository

import (
	"context"

	"budget-book-go/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByProviderID(ctx context.Context, provider string, providerID string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) (*entity.User, error)
}