package repository

import (
	"context"
	"time"

	"budget-book-go/internal/domain/entity"

	"github.com/google/uuid"
)

type IncomeRepository interface {
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Income, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, limit int32, offset int32) ([]*entity.Income, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	FindPlanned(ctx context.Context, userID uuid.UUID) ([]*entity.Income, error)
	FindByDateRange(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]*entity.Income, error)
	Save(ctx context.Context, income *entity.Income) (*entity.Income, error)
	Update(ctx context.Context, income *entity.Income) (*entity.Income, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}