package postgres

import (
	"context"

	"budget-book-go/internal/domain/entity"
	"budget-book-go/internal/domain/repository"
	dbsqlc "budget-book-go/internal/infrastructure/persistence/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db      *pgxpool.Pool
	queries *dbsqlc.Queries
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{
		db:      db,
		queries: dbsqlc.New(db),
	}
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	row, err := r.queries.GetUserByID(ctx, uuidToPgtype(id))
	if err != nil {
		return nil, err
	}
	return userToEntity(row), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return userToEntity(row), nil
}

func (r *userRepository) FindByProviderID(ctx context.Context, provider string, providerID string) (*entity.User, error) {
	row, err := r.queries.GetUserByProviderID(ctx, dbsqlc.GetUserByProviderIDParams{
		Provider:   provider,
		ProviderID: &providerID,
	})
	if err != nil {
		return nil, err
	}
	return userToEntity(row), nil
}

func (r *userRepository) Save(ctx context.Context, user *entity.User) (*entity.User, error) {
	row, err := r.queries.CreateUser(ctx, dbsqlc.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Provider:     user.Provider,
		ProviderID:   user.ProviderID,
		Name:         user.Name,
	})
	if err != nil {
		return nil, err
	}
	return userToEntity(row), nil
}

func userToEntity(row dbsqlc.User) *entity.User {
	return &entity.User{
		ID:           pgtypeToUUID(row.ID),
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Provider:     row.Provider,
		ProviderID:   row.ProviderID,
		Name:         row.Name,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}