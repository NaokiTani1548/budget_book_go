package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("IS_MOCK") == "True" {
			// X-User-Id があればそれを使う
			userIDStr := c.GetHeader("X-User-Id")
			if userIDStr != "" {
				userID, err := uuid.Parse(userIDStr)
				if err == nil {
					c.Set("userID", userID)
					c.Next()
					return
				}
			}
			// X-User-Id がなければデフォルトユーザーを使う
			defaultID := os.Getenv("MOCK_USER_ID")
			if defaultID != "" {
				userID, err := uuid.Parse(defaultID)
				if err == nil {
					c.Set("userID", userID)
					c.Next()
					return
				}
			}
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証トークンが必要です"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証トークンの形式が不正です（Bearer <token>）"})
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証トークンが無効です"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "トークンの解析に失敗しました"})
			return
		}

		userIDStr, ok := claims["userId"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "トークンにユーザーIDが含まれていません"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "ユーザーIDの形式が不正です"})
			return
		}

		// コンテキストにユーザーIDを保存
		c.Set("userID", userID)
		c.Next()
	}
}