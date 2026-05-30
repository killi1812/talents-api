package middleware

import (
	"net/http"
	"os"
	"talents-api/repository"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(apiKeyRepo *repository.APIKeyRepository) gin.HandlerFunc {
	adminKey := os.Getenv("ADMIN_API_KEY")

	return func(c *gin.Context) {
		key := c.GetHeader("X-API-KEY")
		if key == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		if adminKey != "" && key == adminKey {
			c.Set("isAdmin", true)
			c.Next()
			return
		}

		exists, err := apiKeyRepo.Exists(c.Request.Context(), key)
		if err != nil || !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// Update last used asynchronously
		go apiKeyRepo.UpdateLastUsed(c.Request.Context(), key)

		c.Set("isAdmin", false)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, _ := c.Get("isAdmin")
		if isAdmin != true {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
