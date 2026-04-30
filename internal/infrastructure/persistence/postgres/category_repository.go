package postgres

import (
	domainerror "budget-book-go/internal/domain/error"
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

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Category, error) {
	row, err := r.queries.GetCategory(ctx, dbsqlc.GetCategoryParams{
		ID:     uuidToPgtype(id),
		UserID: uuidToPgtype(userID),
	})
	if err != nil {
		return nil, domainerror.NewNotFoundError("category")
	}
	return rowToCategory(row), nil
}

func (r *categoryRepository) Save(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	row, err := r.queries.CreateCategory(ctx, dbsqlc.CreateCategoryParams{
		UserID:    optionalUuidToPgtype(category.UserID),
		Name:      category.Name,
		Type:      category.Type,
		Color:     category.Color,
		SortOrder: int32(category.SortOrder),
		IsDefault: category.IsDefault,
	})
	if err != nil {
		return nil, err
	}
	return rowToCategory(row), nil
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	row, err := r.queries.UpdateCategory(ctx, dbsqlc.UpdateCategoryParams{
		ID:        uuidToPgtype(category.ID),
		UserID:    optionalUuidToPgtype(category.UserID),
		Name:      category.Name,
		Type:      category.Type,
		Color:     category.Color,
		SortOrder: int32(category.SortOrder),
	})
	if err != nil {
		return nil, domainerror.NewNotFoundError("category")
	}
	return rowToCategory(row), nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.queries.DeleteCategory(ctx, dbsqlc.DeleteCategoryParams{
		ID:     uuidToPgtype(id),
		UserID: uuidToPgtype(userID),
	})
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