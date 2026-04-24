package response

import (
	"time"

	"budget-book-go/internal/application/dto"

	"github.com/google/uuid"
)

type ExpenseResponse struct {
	ID            uuid.UUID  `json:"id"`
	Amount        float64    `json:"amount"`
	ExpenseDate   string     `json:"expenseDate"`
	CategoryID    *uuid.UUID `json:"categoryId"`
	CategoryName  *string    `json:"categoryName"`
	Description   *string    `json:"description"`
	PaymentMethod *string    `json:"paymentMethod"`
	Memo          *string    `json:"memo"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func NewExpenseResponse(result *dto.ExpenseResult) ExpenseResponse {
	return ExpenseResponse{
		ID:            result.ID,
		Amount:        result.Amount,
		ExpenseDate:   result.ExpenseDate.Format("2006-01-02"),
		CategoryID:    result.CategoryID,
		CategoryName:  result.CategoryName,
		Description:   result.Description,
		PaymentMethod: result.PaymentMethod,
		Memo:          result.Memo,
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
	}
}

func NewExpenseListResponse(results []*dto.ExpenseResult) []ExpenseResponse {
	responses := make([]ExpenseResponse, len(results))
	for i, result := range results {
		responses[i] = NewExpenseResponse(result)
	}
	return responses
}