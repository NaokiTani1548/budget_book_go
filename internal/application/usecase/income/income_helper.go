package income

import (
	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/entity"
)

func toIncomeResult(i *entity.Income) *dto.IncomeResult {
	return &dto.IncomeResult{
		ID:           i.ID,
		UserID:       i.UserID,
		CategoryID:   i.CategoryID,
		CategoryName: i.CategoryName,
		Amount:       i.Amount,
		Description:  i.Description,
		IncomeDate:   i.IncomeDate,
		Memo:         i.Memo,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
	}
}