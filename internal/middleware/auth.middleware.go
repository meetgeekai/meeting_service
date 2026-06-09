package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.GetHeader("Authorization")
		if value != secret {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid or missing header",
			})
			return
		}
		c.Next()
	}
}
