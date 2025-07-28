package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	authToken "github.com/tph-kds/chat_realtime/backend/internal/api/handlers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Bear Token format
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("Token after trimming Bearer: %s", authHeader) // Log after trimming

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}
		claims, err := authToken.ValidateToken(authHeader)

		if err != nil {
			log.Printf("Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}

}
