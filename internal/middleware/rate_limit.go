package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ginko97/fintech-playground/internal/infrastructure"
)

// RateLimiter simple Redis-based rate limiter
func RateLimit(redis *infrastructure.RedisClient, requests int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "rate:" + c.ClientIP()

		count, _ := redis.GetClient().Incr(c.Request.Context(), key).Result()
		if count == 1 {
			redis.GetClient().Expire(c.Request.Context(), key, duration)
		}

		if count > int64(requests) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
