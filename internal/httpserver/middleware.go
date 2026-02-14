package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
			c.Abort()
			return
		}
		userID, err := h.app.JWTManager.ParseUserID(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
			c.Abort()
			return
		}

		c.Set("userID", userID)

		c.Next()
	}
}
