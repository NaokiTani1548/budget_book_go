package expense

import (
	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/entity"
)

func toExpenseResult(e *entity.Expense) *dto.ExpenseResult {
	return &dto.ExpenseResult{
		ID:            e.ID,
		UserID:        e.UserID,
		CategoryID:    e.CategoryID,
		CategoryName:  e.CategoryName,
		Amount:        e.Amount,
		Description:   e.Description,
		ExpenseDate:   e.ExpenseDate,
		PaymentMethod: e.PaymentMethod,
		Memo:          e.Memo,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}