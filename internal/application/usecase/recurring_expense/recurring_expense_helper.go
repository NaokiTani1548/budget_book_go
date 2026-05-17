package recurringexpense

import (
	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/entity"
)

func toRecurringExpenseResult(re *entity.RecurringExpense) *dto.RecurringExpenseResult {
	return &dto.RecurringExpenseResult{
		ID:            re.ID,
		UserID:        re.UserID,
		CategoryID:    re.CategoryID,
		CategoryName:  re.CategoryName,
		Amount:        re.Amount,
		Description:   re.Description,
		PaymentMethod: re.PaymentMethod,
		Memo:          re.Memo,
		BillingDay:    re.BillingDay,
		StartDate:     re.StartDate,
		EndDate:       re.EndDate,
		IsActive:      re.IsActive,
		CreatedAt:     re.CreatedAt,
		UpdatedAt:     re.UpdatedAt,
	}
}