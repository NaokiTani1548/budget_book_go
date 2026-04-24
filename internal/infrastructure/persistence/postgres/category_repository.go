package postgres

import (
	"context"

	"budget-book-go/internal/domain/entity"
	"budget-book-go/internal/domain/repository"
	dbsqlc "budget-book-go/internal/infrastructure/persistence/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type categoryRepository struct {
	db      *pgxpool.Pool
	queries *dbsqlc.Queries
}

func NewCategoryRepository(db *pgxpool.Pool) repository.CategoryRepository {
	return &categoryRepository{
		db:      db,
		queries: dbsqlc.New(db),
	}
}

func (r *categoryRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error) {
	rows, err := r.queries.ListCategories(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, err
	}

	categories := make([]*entity.Category, len(rows))
	for i, row := range rows {
		categories[i] = rowToCategory(row)
	}
	return categories, nil
}

func rowToCategory(row dbsqlc.Category) *entity.Category {
	return &entity.Category{
		ID:        pgtypeToUUID(row.ID),
		UserID:    optionalPgtypeToUUID(row.UserID),
		Name:      row.Name,
		Type:      row.Type,
		Color:     row.Color,
		SortOrder: int(row.SortOrder),
		IsDefault: row.IsDefault,
		CreatedAt: row.CreatedAt.Time,
	}
}