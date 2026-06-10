package handler

import (
	"io"
	"net/http"
	"strings"
	"log"

	usecaseocr "budget-book-go/internal/application/usecase/ocr"

	"github.com/gin-gonic/gin"
)

type OCRHandler struct {
	analyzeUC *usecaseocr.AnalyzeReceiptUseCase
}

func NewOCRHandler(analyzeUC *usecaseocr.AnalyzeReceiptUseCase) *OCRHandler {
	return &OCRHandler{analyzeUC: analyzeUC}
}

// POST /api/ocr/analyze
func (h *OCRHandler) Analyze(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "画像ファイルが必要です"})
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "画像の読み込みに失敗しました"})
		return
	}

    name := header.Filename
    mimeType := header.Header.Get("Content-Type")
    if mimeType == "" || mimeType == "application/octet-stream" {
        // ファイル名の拡張子から判定
        name := header.Filename
        if strings.HasSuffix(name, ".png") {
            mimeType = "image/png"
        } else if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") {
            mimeType = "image/jpeg"
        } else if strings.HasSuffix(name, ".webp") {
            mimeType = "image/webp"
        } else {
             mimeType = "image/jpeg"
        }
    }
	log.Printf("filename: %s, mimeType: %s", name, mimeType)

	result, err := h.analyzeUC.Execute(c.Request.Context(), imageData, mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}