package request

type CreateIncomeRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	IncomeDate  string  `json:"incomeDate" binding:"required"`
	CategoryID  *string `json:"categoryId"`
	Description *string `json:"description"`
	Memo        *string `json:"memo"`
}

type UpdateIncomeRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	IncomeDate  string  `json:"incomeDate" binding:"required"`
	CategoryID  *string `json:"categoryId"`
	Description *string `json:"description"`
	Memo        *string `json:"memo"`
}