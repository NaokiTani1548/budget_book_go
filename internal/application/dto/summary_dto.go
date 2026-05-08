package dto

import "time"

type ForecastQuery struct {
	UserID     string
	TargetDate time.Time
}

type ForecastResult struct {
	CurrentBalance  float64 `json:"currentBalance"`
	PlannedIncome   float64 `json:"plannedIncome"`
	PlannedExpense  float64 `json:"plannedExpense"`
	ForecastBalance float64 `json:"forecastBalance"`
	TargetDate      string  `json:"targetDate"`
}