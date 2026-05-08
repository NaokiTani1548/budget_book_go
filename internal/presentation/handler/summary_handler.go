package handler

import (
	"net/http"
	"time"

	usecasesummary "budget-book-go/internal/application/usecase/summary"

	"github.com/gin-gonic/gin"
)

type SummaryHandler struct {
	forecastUC *usecasesummary.GetForecastUseCase
}

func NewSummaryHandler(forecastUC *usecasesummary.GetForecastUseCase) *SummaryHandler {
	return &SummaryHandler{forecastUC: forecastUC}
}

// GET /api/summary/forecast?targetDate=2026-06-30
func (h *SummaryHandler) GetForecast(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	targetDate, err := time.Parse("2006-01-02", c.Query("targetDate"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "targetDateの形式が不正です（例: 2026-06-30）"})
		return
	}

	result, err := h.forecastUC.Execute(c.Request.Context(), userID, targetDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}