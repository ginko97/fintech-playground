package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ginko97/fintech-playground/internal/infrastructure"
	"go.uber.org/zap"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // process request

		durationMs := time.Since(start).Milliseconds()

		infrastructure.GetLogger().Info("http_request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("duration_ms", durationMs), // numeric → good for aggregation
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("request_id", getOrGenerateRequestID(c)),
			zap.String("timestamp", time.Now().UTC().Format(time.RFC3339)), // ISO format
		)
	}
}

// Helper function
func getOrGenerateRequestID(c *gin.Context) string {
	if id := c.GetHeader("X-Request-ID"); id != "" {
		return id
	}
	return "req_" + time.Now().Format("20060102150405")
}
