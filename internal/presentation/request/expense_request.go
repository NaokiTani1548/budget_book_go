package request

type CreateExpenseRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	ExpenseDate   string  `json:"expenseDate" binding:"required"`
	CategoryID    *string `json:"categoryId"`
	Description   *string `json:"description"`
	PaymentMethod *string `json:"paymentMethod"`
	Memo          *string `json:"memo"`
	IsPlanned     bool    `json:"isPlanned"`
	PlannedDate   *string `json:"plannedDate"`
}

type UpdateExpenseRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	ExpenseDate   string  `json:"expenseDate" binding:"required"`
	CategoryID    *string `json:"categoryId"`
	Description   *string `json:"description"`
	PaymentMethod *string `json:"paymentMethod"`
	Memo          *string `json:"memo"`
	IsPlanned     bool    `json:"isPlanned"`
	PlannedDate   *string `json:"plannedDate"`
}