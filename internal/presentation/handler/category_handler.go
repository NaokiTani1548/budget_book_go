package handler

import (
	"net/http"

	usecasecategory "budget-book-go/internal/application/usecase/category"
	"budget-book-go/internal/presentation/response"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	getUC *usecasecategory.GetCategoryUseCase
}

func NewCategoryHandler(getUC *usecasecategory.GetCategoryUseCase) *CategoryHandler {
	return &CategoryHandler{getUC: getUC}
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