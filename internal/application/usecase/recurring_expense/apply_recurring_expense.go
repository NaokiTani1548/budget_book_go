package recurringexpense

import (
	"context"
	"time"

	"budget-book-go/internal/domain/entity"
	"budget-book-go/internal/domain/repository"

	"github.com/google/uuid"
)

type ApplyRecurringExpenseUseCase struct {
	recurringRepo repository.RecurringExpenseRepository
	expenseRepo   repository.ExpenseRepository
}

func NewApplyRecurringExpenseUseCase(
	recurringRepo repository.RecurringExpenseRepository,
	expenseRepo repository.ExpenseRepository,
) *ApplyRecurringExpenseUseCase {
	return &ApplyRecurringExpenseUseCase{
		recurringRepo: recurringRepo,
		expenseRepo:   expenseRepo,
	}
}

func (uc *ApplyRecurringExpenseUseCase) Execute(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	// 今月 + 1ヶ月先まで生成
	return uc.applyUpTo(ctx, userID, firstDayOf(now.AddDate(0, 1, 0)))
}

// ExecuteForRange は from〜to の期間に含まれる月の定期支出を適用する。
// 期間指定取得など未来日付を含むクエリ実行前に呼び出す。
func (uc *ApplyRecurringExpenseUseCase) ExecuteForRange(ctx context.Context, userID uuid.UUID, from, to time.Time) error {
	return uc.applyUpTo(ctx, userID, firstDayOf(to))
}

func (uc *ApplyRecurringExpenseUseCase) applyUpTo(ctx context.Context, userID uuid.UUID, upTo time.Time) error {
	list, err := uc.recurringRepo.FindActive(ctx, userID)
	if err != nil {
		return err
	}

	for _, re := range list {
		current := firstDayOf(re.StartDate)

		// EndDate がある場合は upTo と小さい方を上限にする
		target := upTo
		if re.EndDate != nil {
			endMonth := firstDayOf(*re.EndDate)
			if endMonth.Before(target) {
				target = endMonth
			}
		}

		for !current.After(target) {
			year  := current.Year()
			month := int(current.Month())

			exists, err := uc.recurringRepo.ExistsLog(ctx, re.ID, year, month)
			if err != nil {
				return err
			}

			if !exists {
				billingDate := billingDateOf(year, month, re.BillingDay)

				expense := &entity.Expense{
					UserID:        re.UserID,
					CategoryID:    re.CategoryID,
					Amount:        re.Amount,
					Description:   re.Description,
					ExpenseDate:   billingDate,
					PaymentMethod: re.PaymentMethod,
					Memo:          re.Memo,
				}

				saved, err := uc.expenseRepo.Save(ctx, expense)
				if err != nil {
					return err
				}

				if err := uc.recurringRepo.SaveLog(ctx, re.ID, saved.ID, year, month); err != nil {
					return err
				}
			}

			current = current.AddDate(0, 1, 0)
		}
	}
	return nil
}

func firstDayOf(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func billingDateOf(year, month, day int) time.Time {
	// 月末を超える場合は月末に丸める（例: 2月31日 → 2月28日）
	t := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
	lastDay := t.Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}