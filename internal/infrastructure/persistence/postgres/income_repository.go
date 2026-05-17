package postgres

import (
	"context"
	"time"

	"budget-book-go/internal/domain/entity"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/domain/repository"
	dbsqlc "budget-book-go/internal/infrastructure/persistence/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type incomeRepository struct {
	db      *pgxpool.Pool
	queries *dbsqlc.Queries
}

func NewIncomeRepository(db *pgxpool.Pool) repository.IncomeRepository {
	return &incomeRepository{
		db:      db,
		queries: dbsqlc.New(db),
	}
}

func (r *incomeRepository) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Income, error) {
	row, err := r.queries.GetIncome(ctx, dbsqlc.GetIncomeParams{
		ID:     uuidToPgtype(id),
		UserID: uuidToPgtype(userID),
	})
	if err != nil {
		return nil, domainerror.NewNotFoundError("income")
	}
	return getIncomeRowToEntity(row), nil
}

func (r *incomeRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Income, error) {
	rows, err := r.queries.ListIncomes(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, err
	}
	incomes := make([]*entity.Income, len(rows))
	for i, row := range rows {
		incomes[i] = listIncomeRowToEntity(row)
	}
	return incomes, nil
}

func (r *incomeRepository) FindPlanned(ctx context.Context, userID uuid.UUID) ([]*entity.Income, error) {
	rows, err := r.queries.ListPlannedIncomes(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, err
	}
	incomes := make([]*entity.Income, len(rows))
	for i, row := range rows {
		incomes[i] = listPlannedIncomeRowToEntity(row)
	}
	return incomes, nil
}

func (r *incomeRepository) Save(ctx context.Context, income *entity.Income) (*entity.Income, error) {
	row, err := r.queries.CreateIncome(ctx, dbsqlc.CreateIncomeParams{
		UserID:      uuidToPgtype(income.UserID),
		CategoryID:  optionalUuidToPgtype(income.CategoryID),
		Amount:      numericFromFloat(income.Amount),
		Description: income.Description,
		IncomeDate:  dateToPgtype(income.IncomeDate),
		Memo:        income.Memo,
	})
	if err != nil {
		return nil, err
	}
	return savedIncomeToEntity(row), nil
}

func (r *incomeRepository) Update(ctx context.Context, income *entity.Income) (*entity.Income, error) {
	row, err := r.queries.UpdateIncome(ctx, dbsqlc.UpdateIncomeParams{
		ID:          uuidToPgtype(income.ID),
		UserID:      uuidToPgtype(income.UserID),
		CategoryID:  optionalUuidToPgtype(income.CategoryID),
		Amount:      numericFromFloat(income.Amount),
		Description: income.Description,
		IncomeDate:  dateToPgtype(income.IncomeDate),
		Memo:        income.Memo,
	})
	if err != nil {
		return nil, domainerror.NewNotFoundError("income")
	}
	return savedIncomeToEntity(row), nil
}

func (r *incomeRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.queries.DeleteIncome(ctx, dbsqlc.DeleteIncomeParams{
		ID:     uuidToPgtype(id),
		UserID: uuidToPgtype(userID),
	})
}

// -------------------- 変換ヘルパー --------------------

func optionalDateToPgtype(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func optionalPgtypeToTime(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	return &d.Time
}

func getIncomeRowToEntity(row dbsqlc.GetIncomeRow) *entity.Income {
	return &entity.Income{
		ID:           pgtypeToUUID(row.ID),
		UserID:       pgtypeToUUID(row.UserID),
		CategoryID:   optionalPgtypeToUUID(row.CategoryID),
		Amount:       numericToFloat(row.Amount),
		Description:  row.Description,
		IncomeDate:   row.IncomeDate.Time,
		Memo:         row.Memo,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		CategoryName: row.CategoryName,
	}
}

func listIncomeRowToEntity(row dbsqlc.ListIncomesRow) *entity.Income {
	return &entity.Income{
		ID:           pgtypeToUUID(row.ID),
		UserID:       pgtypeToUUID(row.UserID),
		CategoryID:   optionalPgtypeToUUID(row.CategoryID),
		Amount:       numericToFloat(row.Amount),
		Description:  row.Description,
		IncomeDate:   row.IncomeDate.Time,
		Memo:         row.Memo,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		CategoryName: row.CategoryName,
	}
}

func listPlannedIncomeRowToEntity(row dbsqlc.ListPlannedIncomesRow) *entity.Income {
	return &entity.Income{
		ID:           pgtypeToUUID(row.ID),
		UserID:       pgtypeToUUID(row.UserID),
		CategoryID:   optionalPgtypeToUUID(row.CategoryID),
		Amount:       numericToFloat(row.Amount),
		Description:  row.Description,
		IncomeDate:   row.IncomeDate.Time,
		Memo:         row.Memo,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		CategoryName: row.CategoryName,
	}
}

func savedIncomeToEntity(row dbsqlc.Income) *entity.Income {
	return &entity.Income{
		ID:          pgtypeToUUID(row.ID),
		UserID:      pgtypeToUUID(row.UserID),
		CategoryID:  optionalPgtypeToUUID(row.CategoryID),
		Amount:      numericToFloat(row.Amount),
		Description: row.Description,
		IncomeDate:  row.IncomeDate.Time,
		Memo:        row.Memo,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
func (r *incomeRepository) FindByDateRange(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]*entity.Income, error) {
	rows, err := r.queries.ListIncomesByDateRange(ctx, dbsqlc.ListIncomesByDateRangeParams{
		UserID:       uuidToPgtype(userID),
		IncomeDate:   dateToPgtype(from),
		IncomeDate_2: dateToPgtype(to),
	})
	if err != nil {
		return nil, err
	}

	incomes := make([]*entity.Income, len(rows))
	for i, row := range rows {
		incomes[i] = dateRangeIncomeRowToEntity(row)
	}
	return incomes, nil
}

func dateRangeIncomeRowToEntity(row dbsqlc.ListIncomesByDateRangeRow) *entity.Income {
	return &entity.Income{
		ID:           pgtypeToUUID(row.ID),
		UserID:       pgtypeToUUID(row.UserID),
		CategoryID:   optionalPgtypeToUUID(row.CategoryID),
		Amount:       numericToFloat(row.Amount),
		Description:  row.Description,
		IncomeDate:   row.IncomeDate.Time,
		Memo:         row.Memo,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		CategoryName: row.CategoryName,
	}
}