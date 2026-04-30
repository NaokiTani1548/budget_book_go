package repository

import (
	"context"

	"budget-book-go/internal/domain/entity"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error)
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Category, error)
	Save(ctx context.Context, category *entity.Category) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) (*entity.Category, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}