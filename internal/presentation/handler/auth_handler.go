package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	usecaseauth "budget-book-go/internal/application/usecase/auth"

	"github.com/gin-gonic/gin"

	"fmt"
    "net/url"
)

type AuthHandler struct {
	googleAuthUC *usecaseauth.GoogleAuthUseCase
}

func NewAuthHandler(googleAuthUC *usecaseauth.GoogleAuthUseCase) *AuthHandler {
	return &AuthHandler{googleAuthUC: googleAuthUC}
}

// GET /api/auth/google/url
func (h *AuthHandler) GetGoogleAuthURL(c *gin.Context) {
	state := generateState()
	url := h.googleAuthUC.GetAuthURL(state)

	c.JSON(http.StatusOK, gin.H{
		"url":   url,
		"state": state,
	})
}

// GET /api/auth/google/callback?code=xxx&state=xxx
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "認証コードがありません"})
		return
	}

	result, err := h.googleAuthUC.Authenticate(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    frontendURL := os.Getenv("FRONTEND_URL")
    if frontendURL == "" {
        frontendURL = "http://localhost:3000"
    }
    redirectURL := fmt.Sprintf("%s/auth/callback?token=%s&userId=%s&email=%s&name=%s",
        frontendURL,
        url.QueryEscape(result.Token),
        url.QueryEscape(result.UserID),
        url.QueryEscape(result.Email),
        url.QueryEscape(result.Name),
    )
	c.Redirect(http.StatusFound, redirectURL)
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}