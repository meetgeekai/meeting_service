package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

func RateLimiterMiddleware(limiter ratelimit.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter.Take()
		c.Next()
	}
}
