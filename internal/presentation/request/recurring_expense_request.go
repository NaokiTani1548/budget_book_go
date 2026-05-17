package request

type CreateRecurringExpenseRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	BillingDay    int     `json:"billingDay" binding:"required,min=1,max=31"`
	StartDate     string  `json:"startDate" binding:"required"`
	CategoryID    *string `json:"categoryId"`
	Description   *string `json:"description"`
	PaymentMethod *string `json:"paymentMethod"`
	Memo          *string `json:"memo"`
	EndDate       *string `json:"endDate"`
}

type UpdateRecurringExpenseRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	BillingDay    int     `json:"billingDay" binding:"required,min=1,max=31"`
	StartDate     string  `json:"startDate" binding:"required"`
	CategoryID    *string `json:"categoryId"`
	Description   *string `json:"description"`
	PaymentMethod *string `json:"paymentMethod"`
	Memo          *string `json:"memo"`
	EndDate       *string `json:"endDate"`
	IsActive      bool    `json:"isActive"`
}