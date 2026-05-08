package handler

import (
	"net/http"
	"time"
	"budget-book-go/internal/application/dto"
	usecaseincome "budget-book-go/internal/application/usecase/income"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/presentation/request"
	"budget-book-go/internal/presentation/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IncomeHandler struct {
	createUC *usecaseincome.CreateIncomeUseCase
	getUC    *usecaseincome.GetIncomeUseCase
	updateUC *usecaseincome.UpdateIncomeUseCase
	deleteUC *usecaseincome.DeleteIncomeUseCase
}

func NewIncomeHandler(
	createUC *usecaseincome.CreateIncomeUseCase,
	getUC *usecaseincome.GetIncomeUseCase,
	updateUC *usecaseincome.UpdateIncomeUseCase,
	deleteUC *usecaseincome.DeleteIncomeUseCase,
) *IncomeHandler {
	return &IncomeHandler{
		createUC: createUC,
		getUC:    getUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

// GET /api/incomes
func (h *IncomeHandler) GetAll(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	results, err := h.getUC.ExecuteGetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewIncomeListResponse(results))
}

// GET /api/incomes/:id
func (h *IncomeHandler) GetByID(c *gin.Context) {
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

	c.JSON(http.StatusOK, response.NewIncomeResponse(result))
}

// GET /api/incomes/planned
func (h *IncomeHandler) GetPlanned(c *gin.Context) {
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

	c.JSON(http.StatusOK, response.NewIncomeListResponse(results))
}

// GET /api/incomes/date?from=2026-04-01&to=2026-04-30
func (h *IncomeHandler) GetByDateRange(c *gin.Context) {
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

	results, err := h.getUC.ExecuteGetByDateRange(c.Request.Context(), userID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewIncomeListResponse(results))
}

// POST /api/incomes
func (h *IncomeHandler) Create(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	var req request.CreateIncomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incomeDate, err := time.Parse("2006-01-02", req.IncomeDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incomeDateの形式が不正です（例: 2026-04-14）"})
		return
	}

	categoryID, err := parseOptionalUUID(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "categoryIdの形式が不正です"})
		return
	}

	plannedDate, err := parseOptionalDate(req.PlannedDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plannedDateの形式が不正です（例: 2026-05-25）"})
		return
	}

	cmd := dto.CreateIncomeCommand{
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      req.Amount,
		Description: req.Description,
		IncomeDate:  incomeDate,
		Memo:        req.Memo,
		IsPlanned:   req.IsPlanned,
		PlannedDate: plannedDate,
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

	c.JSON(http.StatusCreated, response.NewIncomeResponse(result))
}

// PUT /api/incomes/:id
func (h *IncomeHandler) Update(c *gin.Context) {
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

	var req request.UpdateIncomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incomeDate, err := time.Parse("2006-01-02", req.IncomeDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incomeDateの形式が不正です（例: 2026-04-14）"})
		return
	}

	categoryID, err := parseOptionalUUID(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "categoryIdの形式が不正です"})
		return
	}

	plannedDate, err := parseOptionalDate(req.PlannedDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plannedDateの形式が不正です（例: 2026-05-25）"})
		return
	}

	cmd := dto.UpdateIncomeCommand{
		ID:          id,
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      req.Amount,
		Description: req.Description,
		IncomeDate:  incomeDate,
		Memo:        req.Memo,
		IsPlanned:   req.IsPlanned,
		PlannedDate: plannedDate,
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

	c.JSON(http.StatusOK, response.NewIncomeResponse(result))
}

// DELETE /api/incomes/:id
func (h *IncomeHandler) Delete(c *gin.Context) {
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