package handler

import (
	recurringexpense "budget-book-go/internal/application/usecase/recurring_expense"
	"net/http"
	"time"
	"fmt"

	"budget-book-go/internal/application/dto"
	usecaseexpense "budget-book-go/internal/application/usecase/expense"
	domainerror "budget-book-go/internal/domain/error"
	repository "budget-book-go/internal/domain/repository"
	"budget-book-go/internal/presentation/request"
	"budget-book-go/internal/presentation/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExpenseHandler struct {
	createUC *usecaseexpense.CreateExpenseUseCase
	getUC    *usecaseexpense.GetExpenseUseCase
	updateUC *usecaseexpense.UpdateExpenseUseCase
	deleteUC *usecaseexpense.DeleteExpenseUseCase
	applyUC  *recurringexpense.ApplyRecurringExpenseUseCase
}

func NewExpenseHandler(
	createUC *usecaseexpense.CreateExpenseUseCase,
	getUC *usecaseexpense.GetExpenseUseCase,
	updateUC *usecaseexpense.UpdateExpenseUseCase,
	deleteUC *usecaseexpense.DeleteExpenseUseCase,
	applyUC *recurringexpense.ApplyRecurringExpenseUseCase,
) *ExpenseHandler {
	return &ExpenseHandler{
		createUC: createUC,
		getUC:    getUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		applyUC:  applyUC,
	}
}

// GET /api/expenses
func (h *ExpenseHandler) GetAllByUserID(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	_ = h.applyUC.Execute(c.Request.Context(), userID)

	results, err := h.getUC.ExecuteGetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewExpenseListResponse(results))
}

// GET /api/expenses/planned
func (h *ExpenseHandler) GetPlanned(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	results, err := h.getUC.ExecuteGetPlanned(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewExpenseListResponse(results))
}

// GET /api/expenses/:id
func (h *ExpenseHandler) GetByID(c *gin.Context) {
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

	c.JSON(http.StatusOK, response.NewExpenseResponse(result))
}

// GET /api/expenses/date?from=2026-04-01&to=2026-04-30
func (h *ExpenseHandler) GetByDateRange(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	from, err := time.Parse("2006-01-02", c.Query("from"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fromの形式が不正です（例: 2026-04-01）"})
		return
	}

	to, err := time.Parse("2006-01-02", c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "toの形式が不正です（例: 2026-04-30）"})
		return
	}

	_ = h.applyUC.ExecuteForRange(c.Request.Context(), userID, from, to)

	results, err := h.getUC.ExecuteGetByDateRange(c.Request.Context(), userID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewExpenseListResponse(results))
}

// POST /api/expenses
func (h *ExpenseHandler) Create(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	var req request.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日付の形式が不正です（例: 2026-04-14）"})
		return
	}

	categoryID, err := parseOptionalUUID(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "categoryIdの形式が不正です"})
		return
	}

	cmd := dto.CreateExpenseCommand{
		UserID:        userID,
		CategoryID:    categoryID,
		Amount:        req.Amount,
		Description:   req.Description,
		ExpenseDate:   expenseDate,
		PaymentMethod: req.PaymentMethod,
		Memo:          req.Memo,
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

	c.JSON(http.StatusCreated, response.NewExpenseResponse(result))
}

// PUT /api/expenses/:id
func (h *ExpenseHandler) Update(c *gin.Context) {
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

	var req request.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日付の形式が不正です（例: 2026-04-14）"})
		return
	}

	categoryID, err := parseOptionalUUID(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "categoryIdの形式が不正です"})
		return
	}

	cmd := dto.UpdateExpenseCommand{
		ID:            id,
		UserID:        userID,
		CategoryID:    categoryID,
		Amount:        req.Amount,
		Description:   req.Description,
		ExpenseDate:   expenseDate,
		PaymentMethod: req.PaymentMethod,
		Memo:          req.Memo,
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

	c.JSON(http.StatusOK, response.NewExpenseResponse(result))
}

// DELETE /api/expenses/:id
func (h *ExpenseHandler) Delete(c *gin.Context) {
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

// -------------------- ヘルパー --------------------

func extractUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, fmt.Errorf("userID not found in context")
	}
	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("userID is not uuid.UUID")
	}
	return id, nil
}

func parseOptionalUUID(s *string) (*uuid.UUID, error) {
	if s == nil {
		return nil, nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func isDomainError(err error, code domainerror.ErrorCode) bool {
	if de, ok := err.(*domainerror.DomainError); ok {
		return de.Code == code
	}
	return false
}

func parseOptionalDate(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GET /api/expenses/search?from=&to=&categoryId=&keyword=
func (h *ExpenseHandler) Search(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "認証エラー"})
		return
	}

	var params repository.SearchExpenseParams

	if from := c.Query("from"); from != "" {
		t, err := time.Parse("2006-01-02", from)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fromの形式が不正です（例: 2026-05-01）"})
			return
		}
		params.DateFrom = &t
	}

	if to := c.Query("to"); to != "" {
		t, err := time.Parse("2006-01-02", to)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "toの形式が不正です（例: 2026-05-31）"})
			return
		}
		params.DateTo = &t
	}

	if categoryID := c.Query("categoryId"); categoryID != "" {
		id, err := uuid.Parse(categoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "categoryIdの形式が不正です"})
			return
		}
		params.CategoryID = &id
	}

	if keyword := c.Query("keyword"); keyword != "" {
		params.Keyword = &keyword
	}

	results, err := h.getUC.ExecuteSearch(c.Request.Context(), userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewExpenseListResponse(results))
}