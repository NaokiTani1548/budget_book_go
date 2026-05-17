package postgres

import (
	"context"
	"time"

	"budget-book-go/internal/domain/entity"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/domain/repository"
	dbsqlc "budget-book-go/internal/infrastructure/persistence/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type recurringExpenseRepository struct {
	db      *pgxpool.Pool
	queries *dbsqlc.Queries
}

func NewRecurringExpenseRepository(db *pgxpool.Pool) repository.RecurringExpenseRepository {
	return &recurringExpenseRepository{
		db:      db,
		queries: dbsqlc.New(db),
	}
}

func (r *recurringExpenseRepository) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.RecurringExpense, error) {
	row, err := r.queries.GetRecurringExpense(ctx, dbsqlc.GetRecurringExpenseParams{
		ID:     uuidToPgtype(id),
		UserID: uuidToPgtype(userID),
	})
	if err != nil {
		return nil, domainerror.NewNotFoundError("recurring_expense")
	}
	return getRecurringExpenseRowToEntity(row), nil
}

func (r *recurringExpenseRepository) FindAll(ctx context.Context, userID uuid.UUID) ([]*entity.RecurringExpense, error) {
	rows, err := r.queries.ListRecurringExpenses(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, err
	}
	return listRecurringExpenseRowsToEntities(rows), nil
}

func (r *recurringExpenseRepository) FindActive(ctx context.Context, userID uuid.UUID) ([]*entity.RecurringExpense, error) {
	rows, err := r.queries.ListActiveRecurringExpenses(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, err
	}
	result := make([]*entity.RecurringExpense, len(rows))
	for i, row := range rows {
		result[i] = listActiveRecurringExpenseRowToEntity(row)
	}
	return result, nil
}

func (r *recurringExpenseRepository) Save(ctx context.Context, re *entity.RecurringExpense) (*entity.RecurringExpense, error) {
	row, err := r.queries.CreateRecurringExpense(ctx, dbsqlc.CreateRecurringExpenseParams{
		UserID:        uuidToPgtype(re.UserID),
		CategoryID:    optionalUuidToPgtype(re.CategoryID),
		Amount:        numericFromFloat(re.Amount),
		Description:   re.Description,
		PaymentMethod: re.PaymentMethod,
		Memo:          re.Memo,
		BillingDay:    int32(re.BillingDay),
		StartDate:     dateToPgtype(re.StartDate),
		EndDate:       optionalDateToPgtype(re.EndDate),
	})
	if err != nil {
		return nil, err
	}
	return recurringExpenseToEntity(row), nil
}

func (r *recurringExpenseRepository) Update(ctx context.Context, re *entity.RecurringExpense) (*entity.RecurringExpense, error) {
	row, err := r.queries.UpdateRecurringExpense(ctx, dbsqlc.UpdateRecurringExpenseParams{
		ID:            uuidToPgtype(re.ID),
		UserID:        uuidToPgtype(re.UserID),
		CategoryID:    optionalUuidToPgtype(re.CategoryID),
		Amount:        numericFromFloat(re.Amount),
		Description:   re.Description,
		PaymentMethod: re.PaymentMethod,
		Memo:          re.Memo,
		BillingDay:    int32(re.BillingDay),
		StartDate:     dateToPgtype(re.StartDate),
		EndDate:       optionalDateToPgtype(re.EndDate),
		IsActive:      re.IsActive,
	})
	if err != nil {
		return nil, domainerror.NewNotFoundError("recurring_expense")
	}
	return recurringExpenseToEntity(row), nil
}

func (r *recurringExpenseRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.queries.DeleteRecurringExpense(ctx, dbsqlc.DeleteRecurringExpenseParams{
		ID:     uuidToPgtype(id),
		UserID: uuidToPgtype(userID),
	})
}

func (r *recurringExpenseRepository) ExistsLog(ctx context.Context, recurringExpenseID uuid.UUID, year int, month int) (bool, error) {
	_, err := r.queries.GetRecurringExpenseLog(ctx, dbsqlc.GetRecurringExpenseLogParams{
		RecurringExpenseID: uuidToPgtype(recurringExpenseID),
		BillingYear:        int32(year),
		BillingMonth:       int32(month),
	})
	if err != nil {
		return false, nil // 見つからない = 未処理
	}
	return true, nil
}

func (r *recurringExpenseRepository) SaveLog(ctx context.Context, recurringExpenseID uuid.UUID, expenseID uuid.UUID, year int, month int) error {
	_, err := r.queries.CreateRecurringExpenseLog(ctx, dbsqlc.CreateRecurringExpenseLogParams{
		RecurringExpenseID: uuidToPgtype(recurringExpenseID),
		ExpenseID:          uuidToPgtype(expenseID),
		BillingYear:        int32(year),
		BillingMonth:       int32(month),
	})
	return err
}

// -------------------- 変換ヘルパー --------------------

func getRecurringExpenseRowToEntity(row dbsqlc.GetRecurringExpenseRow) *entity.RecurringExpense {
	return &entity.RecurringExpense{
		ID:            pgtypeToUUID(row.ID),
		UserID:        pgtypeToUUID(row.UserID),
		CategoryID:    optionalPgtypeToUUID(row.CategoryID),
		Amount:        numericToFloat(row.Amount),
		Description:   row.Description,
		PaymentMethod: row.PaymentMethod,
		Memo:          row.Memo,
		BillingDay:    int(row.BillingDay),
		StartDate:     row.StartDate.Time,
		EndDate:       optionalPgtypeToTime(row.EndDate),
		IsActive:      row.IsActive,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
		CategoryName:  row.CategoryName,
	}
}

func listRecurringExpenseRowsToEntities(rows []dbsqlc.ListRecurringExpensesRow) []*entity.RecurringExpense {
	result := make([]*entity.RecurringExpense, len(rows))
	for i, row := range rows {
		result[i] = &entity.RecurringExpense{
			ID:            pgtypeToUUID(row.ID),
			UserID:        pgtypeToUUID(row.UserID),
			CategoryID:    optionalPgtypeToUUID(row.CategoryID),
			Amount:        numericToFloat(row.Amount),
			Description:   row.Description,
			PaymentMethod: row.PaymentMethod,
			Memo:          row.Memo,
			BillingDay:    int(row.BillingDay),
			StartDate:     row.StartDate.Time,
			EndDate:       optionalPgtypeToTime(row.EndDate),
			IsActive:      row.IsActive,
			CreatedAt:     row.CreatedAt.Time,
			UpdatedAt:     row.UpdatedAt.Time,
			CategoryName:  row.CategoryName,
		}
	}
	return result
}

func listActiveRecurringExpenseRowToEntity(row dbsqlc.ListActiveRecurringExpensesRow) *entity.RecurringExpense {
	return &entity.RecurringExpense{
		ID:            pgtypeToUUID(row.ID),
		UserID:        pgtypeToUUID(row.UserID),
		CategoryID:    optionalPgtypeToUUID(row.CategoryID),
		Amount:        numericToFloat(row.Amount),
		Description:   row.Description,
		PaymentMethod: row.PaymentMethod,
		Memo:          row.Memo,
		BillingDay:    int(row.BillingDay),
		StartDate:     row.StartDate.Time,
		EndDate:       optionalPgtypeToTime(row.EndDate),
		IsActive:      row.IsActive,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
		CategoryName:  row.CategoryName,
	}
}

func recurringExpenseToEntity(row dbsqlc.RecurringExpense) *entity.RecurringExpense {
	return &entity.RecurringExpense{
		ID:            pgtypeToUUID(row.ID),
		UserID:        pgtypeToUUID(row.UserID),
		CategoryID:    optionalPgtypeToUUID(row.CategoryID),
		Amount:        numericToFloat(row.Amount),
		Description:   row.Description,
		PaymentMethod: row.PaymentMethod,
		Memo:          row.Memo,
		BillingDay:    int(row.BillingDay),
		StartDate:     row.StartDate.Time,
		EndDate:       optionalPgtypeToTime(row.EndDate),
		IsActive:      row.IsActive,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}

// optionalDateToPgtype は income_repository.go に定義済みのため不要
// optionalPgtypeToTime は income_repository.go に定義済みのため不要
var _ time.Time // time パッケージの使用を明示