package response

import (
	"time"

	"budget-book-go/internal/application/dto"

	"github.com/google/uuid"
)

type IncomeResponse struct {
	ID           uuid.UUID  `json:"id"`
	Amount       float64    `json:"amount"`
	IncomeDate   string     `json:"incomeDate"`
	CategoryID   *uuid.UUID `json:"categoryId"`
	CategoryName *string    `json:"categoryName"`
	Description  *string    `json:"description"`
	Memo         *string    `json:"memo"`
	IsPlanned    bool       `json:"isPlanned"`
	PlannedDate  *string    `json:"plannedDate"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

func NewIncomeResponse(result *dto.IncomeResult) IncomeResponse {
	var plannedDate *string
	if result.PlannedDate != nil {
		s := result.PlannedDate.Format("2006-01-02")
		plannedDate = &s
	}

	return IncomeResponse{
		ID:           result.ID,
		Amount:       result.Amount,
		IncomeDate:   result.IncomeDate.Format("2006-01-02"),
		CategoryID:   result.CategoryID,
		CategoryName: result.CategoryName,
		Description:  result.Description,
		Memo:         result.Memo,
		IsPlanned:    result.IsPlanned,
		PlannedDate:  plannedDate,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}
}

func NewIncomeListResponse(results []*dto.IncomeResult) []IncomeResponse {
	responses := make([]IncomeResponse, len(results))
	for i, result := range results {
		responses[i] = NewIncomeResponse(result)
	}
	return responses
}