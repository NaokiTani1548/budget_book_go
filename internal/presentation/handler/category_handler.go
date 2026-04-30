package handler

import (
	"budget-book-go/internal/application/dto"
	domainerror "budget-book-go/internal/domain/error"
	"budget-book-go/internal/presentation/request"
	"net/http"

	usecasecategory "budget-book-go/internal/application/usecase/category"
	"budget-book-go/internal/presentation/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	getUC *usecasecategory.GetCategoryUseCase
	createUC *usecasecategory.CreateCategoryUseCase
	updateUC *usecasecategory.UpdateCategoryUseCase
	deleteUC *usecasecategory.DeleteCategoryUseCase
}

func NewCategoryHandler(getUC *usecasecategory.GetCategoryUseCase,
	createUC *usecasecategory.CreateCategoryUseCase,
	updateUC *usecasecategory.UpdateCategoryUseCase,
	deleteUC *usecasecategory.DeleteCategoryUseCase,
	) *CategoryHandler {
	return &CategoryHandler{
		getUC:    getUC,
		createUC: createUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

// GET /api/categories
func (h *CategoryHandler) GetAllByUserID(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	results, err := h.getUC.ExecuteGetAllByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewCategoryListResponse(results))
}

// POST /api/categories
func (h *CategoryHandler) Create(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Idヘッダーが不正です"})
		return
	}

	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := dto.CreateCategoryCommand{
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Color:     req.Color,
		SortOrder: req.SortOrder,
		IsDefault: req.IsDefault,
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

	c.JSON(http.StatusCreated, response.NewCategoryResponse(result))
}

// PUT /api/categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
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

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := dto.UpdateCategoryCommand{
		ID:        id,
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Color:     req.Color,
		SortOrder: req.SortOrder,
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

	c.JSON(http.StatusOK, response.NewCategoryResponse(result))
}

// DELETE /api/categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
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