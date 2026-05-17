package recurringexpense

import (
	"context"

	"budget-book-go/internal/application/dto"
	"budget-book-go/internal/domain/entity"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/domain/repository"
	"budget-book-go/internal/domain/valueobject"
)

type CreateRecurringExpenseUseCase struct {
	recurringRepo repository.RecurringExpenseRepository
}

func NewCreateRecurringExpenseUseCase(recurringRepo repository.RecurringExpenseRepository) *CreateRecurringExpenseUseCase {
	return &CreateRecurringExpenseUseCase{recurringRepo: recurringRepo}
}

func (uc *CreateRecurringExpenseUseCase) Execute(ctx context.Context, cmd dto.CreateRecurringExpenseCommand) (*dto.RecurringExpenseResult, error) {
	_, err := valueobject.NewMoney(cmd.Amount)
	if err != nil {
		return nil, domainerror.NewInvalidInputError(err.Error())
	}

	if cmd.BillingDay < 1 || cmd.BillingDay > 31 {
		return nil, domainerror.NewInvalidInputError("billing_dayは1〜31の範囲で指定してください")
	}

	re := &entity.RecurringExpense{
		UserID:        cmd.UserID,
		CategoryID:    cmd.CategoryID,
		Amount:        cmd.Amount,
		Description:   cmd.Description,
		PaymentMethod: cmd.PaymentMethod,
		Memo:          cmd.Memo,
		BillingDay:    cmd.BillingDay,
		StartDate:     cmd.StartDate,
		EndDate:       cmd.EndDate,
		IsActive:      true,
	}

	saved, err := uc.recurringRepo.Save(ctx, re)
	if err != nil {
		return nil, err
	}
	return toRecurringExpenseResult(saved), nil
}