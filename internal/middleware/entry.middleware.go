package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const CORRELATION_ID_KEY = "correlation_id"

func EntryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(CORRELATION_ID_KEY, uuid.New().String())
		c.Next()
	}
}
