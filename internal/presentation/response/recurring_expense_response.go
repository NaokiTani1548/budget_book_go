package response

import (
	"time"

	"budget-book-go/internal/application/dto"

	"github.com/google/uuid"
)

type RecurringExpenseResponse struct {
	ID            uuid.UUID  `json:"id"`
	Amount        float64    `json:"amount"`
	BillingDay    int        `json:"billingDay"`
	StartDate     string     `json:"startDate"`
	EndDate       *string    `json:"endDate"`
	CategoryID    *uuid.UUID `json:"categoryId"`
	CategoryName  *string    `json:"categoryName"`
	Description   *string    `json:"description"`
	PaymentMethod *string    `json:"paymentMethod"`
	Memo          *string    `json:"memo"`
	IsActive      bool       `json:"isActive"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func NewRecurringExpenseResponse(result *dto.RecurringExpenseResult) RecurringExpenseResponse {
	var endDate *string
	if result.EndDate != nil {
		s := result.EndDate.Format("2006-01-02")
		endDate = &s
	}

	return RecurringExpenseResponse{
		ID:            result.ID,
		Amount:        result.Amount,
		BillingDay:    result.BillingDay,
		StartDate:     result.StartDate.Format("2006-01-02"),
		EndDate:       endDate,
		CategoryID:    result.CategoryID,
		CategoryName:  result.CategoryName,
		Description:   result.Description,
		PaymentMethod: result.PaymentMethod,
		Memo:          result.Memo,
		IsActive:      result.IsActive,
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
	}
}

func NewRecurringExpenseListResponse(results []*dto.RecurringExpenseResult) []RecurringExpenseResponse {
	responses := make([]RecurringExpenseResponse, len(results))
	for i, result := range results {
		responses[i] = NewRecurringExpenseResponse(result)
	}
	return responses
}