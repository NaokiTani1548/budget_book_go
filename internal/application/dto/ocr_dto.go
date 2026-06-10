package dto

type OCRItem struct {
	Description   string  `json:"description"`
	Amount        float64 `json:"amount"`
	ExpenseDate   string  `json:"expenseDate"`
	PaymentMethod string  `json:"paymentMethod"`
}

type OCRResult struct {
	Items []OCRItem `json:"items"`
}