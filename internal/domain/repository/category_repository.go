package repository

import (
	"context"

	"budget-book-go/internal/domain/entity"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error)
}