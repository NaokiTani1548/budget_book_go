package handler

import (
	"net/http"
	"time"

	"budget-book-go/internal/application/dto"
	recurringexpense "budget-book-go/internal/application/usecase/recurring_expense"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/presentation/request"
	"budget-book-go/internal/presentation/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecurringExpenseHandler struct {
	getUC    *recurringexpense.GetRecurringExpenseUseCase
	createUC *recurringexpense.CreateRecurringExpenseUseCase
	updateUC *recurringexpense.UpdateRecurringExpenseUseCase
	deleteUC *recurringexpense.DeleteRecurringExpenseUseCase
	applyUC  *recurringexpense.ApplyRecurringExpenseUseCase
}

func NewRecurringExpenseHandler(
	getUC *recurringexpense.GetRecurringExpenseUseCase,
	createUC *recurringexpense.CreateRecurringExpenseUseCase,
	updateUC *recurringexpense.UpdateRecurringExpenseUseCase,
	deleteUC *recurringexpense.DeleteRecurringExpenseUseCase,
	applyUC *recurringexpense.ApplyRecurringExpenseUseCase,
) *RecurringExpenseHandler {
	return &RecurringExpenseHandler{
		getUC:    getUC,
		createUC: createUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		applyUC:  applyUC,
	}
}

// GET /api/recurring-expenses
func (h *RecurringExpenseHandler) GetAll(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	// APIリクエスト時に未処理の定期支出を自動生成
	_ = h.applyUC.Execute(c.Request.Context(), userID)

	results, err := h.getUC.ExecuteGetAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewRecurringExpenseListResponse(results))
}

// GET /api/recurring-expenses/:id
func (h *RecurringExpenseHandler) GetByID(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDの形式が不正です"})
		return
	}

	result, err := h.getUC.ExecuteGetOne(c.Request.Context(), id, userID)
	if err != nil {
		if isDomainError(err, domainerror.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewRecurringExpenseResponse(result))
}

// POST /api/recurring-expenses
func (h *RecurringExpenseHandler) Create(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	var req request.CreateRecurringExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "startDateの形式が不正です（例: 2026-04-01）"})
		return
	}

	endDate, err := parseOptionalDate(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "endDateの形式が不正です（例: 2027-03-31）"})
		return
	}

	categoryID, err := parseOptionalUUID(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "categoryIdの形式が不正です"})
		return
	}

	cmd := dto.CreateRecurringExpenseCommand{
		UserID:        userID,
		CategoryID:    categoryID,
		Amount:        req.Amount,
		Description:   req.Description,
		PaymentMethod: req.PaymentMethod,
		Memo:          req.Memo,
		BillingDay:    req.BillingDay,
		StartDate:     startDate,
		EndDate:       endDate,
	}

	result, err := h.createUC.Execute(c.Request.Context(), cmd)
	if err != nil {
		if isDomainError(err, domainerror.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.NewRecurringExpenseResponse(result))
}

// PUT /api/recurring-expenses/:id
func (h *RecurringExpenseHandler) Update(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDの形式が不正です"})
		return
	}

	var req request.UpdateRecurringExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "startDateの形式が不正です（例: 2026-04-01）"})
		return
	}

	endDate, err := parseOptionalDate(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "endDateの形式が不正です（例: 2027-03-31）"})
		return
	}

	categoryID, err := parseOptionalUUID(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "categoryIdの形式が不正です"})
		return
	}

	cmd := dto.UpdateRecurringExpenseCommand{
		ID:            id,
		UserID:        userID,
		CategoryID:    categoryID,
		Amount:        req.Amount,
		Description:   req.Description,
		PaymentMethod: req.PaymentMethod,
		Memo:          req.Memo,
		BillingDay:    req.BillingDay,
		StartDate:     startDate,
		EndDate:       endDate,
		IsActive:      req.IsActive,
	}

	result, err := h.updateUC.Execute(c.Request.Context(), cmd)
	if err != nil {
		if isDomainError(err, domainerror.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if isDomainError(err, domainerror.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewRecurringExpenseResponse(result))
}

// DELETE /api/recurring-expenses/:id
func (h *RecurringExpenseHandler) Delete(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDの形式が不正です"})
		return
	}

	if err := h.deleteUC.Execute(c.Request.Context(), id, userID); err != nil {
		if isDomainError(err, domainerror.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}