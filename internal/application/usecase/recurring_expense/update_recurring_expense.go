package recurringexpense

import (
	"context"

	"budget-book-go/internal/application/dto"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/domain/repository"
	"budget-book-go/internal/domain/valueobject"
)

type UpdateRecurringExpenseUseCase struct {
	recurringRepo repository.RecurringExpenseRepository
}

func NewUpdateRecurringExpenseUseCase(recurringRepo repository.RecurringExpenseRepository) *UpdateRecurringExpenseUseCase {
	return &UpdateRecurringExpenseUseCase{recurringRepo: recurringRepo}
}

func (uc *UpdateRecurringExpenseUseCase) Execute(ctx context.Context, cmd dto.UpdateRecurringExpenseCommand) (*dto.RecurringExpenseResult, error) {
	_, err := valueobject.NewMoney(cmd.Amount)
	if err != nil {
		return nil, domainerror.NewInvalidInputError(err.Error())
	}

	if cmd.BillingDay < 1 || cmd.BillingDay > 31 {
		return nil, domainerror.NewInvalidInputError("billing_dayは1〜31の範囲で指定してください")
	}

	re, err := uc.recurringRepo.FindByID(ctx, cmd.ID, cmd.UserID)
	if err != nil {
		return nil, err
	}

	re.CategoryID    = cmd.CategoryID
	re.Amount        = cmd.Amount
	re.Description   = cmd.Description
	re.PaymentMethod = cmd.PaymentMethod
	re.Memo          = cmd.Memo
	re.BillingDay    = cmd.BillingDay
	re.StartDate     = cmd.StartDate
	re.EndDate       = cmd.EndDate
	re.IsActive      = cmd.IsActive

	updated, err := uc.recurringRepo.Update(ctx, re)
	if err != nil {
		return nil, err
	}
	return toRecurringExpenseResult(updated), nil
}