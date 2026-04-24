package postgres

import (
	"context"
	"fmt"
	"time"

	"budget-book-go/internal/domain/entity"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/domain/repository"
	dbsqlc "budget-book-go/internal/infrastructure/persistence/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type expenseRepository struct {
	db      *pgxpool.Pool
	queries *dbsqlc.Queries
}

func NewExpenseRepository(db *pgxpool.Pool) repository.ExpenseRepository {
	return &expenseRepository{
		db:      db,
		queries: dbsqlc.New(db),
	}
}

// -------------------- 公開メソッド --------------------

func (r *expenseRepository) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Expense, error) {
	row, err := r.queries.GetExpense(ctx, dbsqlc.GetExpenseParams{
		ID:     uuidToPgtype(id),
		UserID: uuidToPgtype(userID),
	})
	if err != nil {
		return nil, domainerror.NewNotFoundError("expense")
	}
	return rowToExpense(row), nil
}

func (r *expenseRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Expense, error) {
	rows, err := r.queries.ListExpenses(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, err
	}

	expenses := make([]*entity.Expense, len(rows))
	for i, row := range rows {
		expenses[i] = listRowToExpense(row)
	}
	return expenses, nil
}

func (r *expenseRepository) FindByDateRange(ctx context.Context, userID uuid.UUID, from time.Time, to time.Time) ([]*entity.Expense, error) {
	// Phase3で実装
	return nil, nil
}

func (r *expenseRepository) Save(ctx context.Context, expense *entity.Expense) (*entity.Expense, error) {
	row, err := r.queries.CreateExpense(ctx, dbsqlc.CreateExpenseParams{
		UserID:        uuidToPgtype(expense.UserID),
		CategoryID:    optionalUuidToPgtype(expense.CategoryID),
		Amount:        numericFromFloat(expense.Amount),
		Description:   expense.Description,
		ExpenseDate:   dateToPgtype(expense.ExpenseDate),
		PaymentMethod: expense.PaymentMethod,
		Memo:          expense.Memo,
	})
	if err != nil {
		return nil, err
	}
	return savedExpenseToEntity(row), nil
}

func (r *expenseRepository) Update(ctx context.Context, expense *entity.Expense) (*entity.Expense, error) {
	row, err := r.queries.UpdateExpense(ctx, dbsqlc.UpdateExpenseParams{
		ID:            uuidToPgtype(expense.ID),
		UserID:        uuidToPgtype(expense.UserID),
		CategoryID:    optionalUuidToPgtype(expense.CategoryID),
		Amount:        numericFromFloat(expense.Amount),
		Description:   expense.Description,
		ExpenseDate:   dateToPgtype(expense.ExpenseDate),
		PaymentMethod: expense.PaymentMethod,
		Memo:          expense.Memo,
	})
	if err != nil {
		return nil, domainerror.NewNotFoundError("expense")
	}
	return savedExpenseToEntity(row), nil
}

func (r *expenseRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.queries.DeleteExpense(ctx, dbsqlc.DeleteExpenseParams{
		ID:     uuidToPgtype(id),
		UserID: uuidToPgtype(userID),
	})
}

// -------------------- 変換ヘルパー --------------------

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func optionalUuidToPgtype(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

func dateToPgtype(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func numericFromFloat(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	if err := n.Scan(fmt.Sprintf("%f", f)); err != nil {
		return pgtype.Numeric{}
	}
	return n
}

func pgtypeToUUID(p pgtype.UUID) uuid.UUID {
	return uuid.UUID(p.Bytes)
}

func optionalPgtypeToUUID(p pgtype.UUID) *uuid.UUID {
	if !p.Valid {
		return nil
	}
	id := uuid.UUID(p.Bytes)
	return &id
}

func numericToFloat(n pgtype.Numeric) float64 {
	f, _ := n.Float64Value()
	return f.Float64
}

func rowToExpense(row dbsqlc.GetExpenseRow) *entity.Expense {
	return &entity.Expense{
		ID:            pgtypeToUUID(row.ID),
		UserID:        pgtypeToUUID(row.UserID),
		CategoryID:    optionalPgtypeToUUID(row.CategoryID),
		Amount:        numericToFloat(row.Amount),
		Description:   row.Description,
		ExpenseDate:   row.ExpenseDate.Time,
		PaymentMethod: row.PaymentMethod,
		Memo:          row.Memo,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
		CategoryName:  row.CategoryName,
	}
}

func listRowToExpense(row dbsqlc.ListExpensesRow) *entity.Expense {
	return &entity.Expense{
		ID:            pgtypeToUUID(row.ID),
		UserID:        pgtypeToUUID(row.UserID),
		CategoryID:    optionalPgtypeToUUID(row.CategoryID),
		Amount:        numericToFloat(row.Amount),
		Description:   row.Description,
		ExpenseDate:   row.ExpenseDate.Time,
		PaymentMethod: row.PaymentMethod,
		Memo:          row.Memo,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
		CategoryName:  row.CategoryName,
	}
}


func savedExpenseToEntity(row dbsqlc.Expense) *entity.Expense {
	return &entity.Expense{
		ID:            pgtypeToUUID(row.ID),
		UserID:        pgtypeToUUID(row.UserID),
		CategoryID:    optionalPgtypeToUUID(row.CategoryID),
		Amount:        numericToFloat(row.Amount),
		Description:   row.Description,
		ExpenseDate:   row.ExpenseDate.Time,
		PaymentMethod: row.PaymentMethod,
		Memo:          row.Memo,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}