package middleware

import (
	"net/http"
	"sync"
	"time"

	"healthcare-api/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burst,
	}
}

// getLimiter gets or creates a limiter for a client
func (rl *RateLimiter) getLimiter(clientID string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[clientID]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check after acquiring write lock
		if limiter, exists = rl.limiters[clientID]; !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.limiters[clientID] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

// RateLimit middleware applies rate limiting per client IP
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := rl.getLimiter(clientIP)

		if !limiter.Allow() {
			c.Header("X-RateLimit-Limit", "100")
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", time.Now().Add(time.Minute).Format(time.RFC3339))
			
			c.JSON(http.StatusTooManyRequests, models.NewOperationOutcome("error", "throttled", "Rate limit exceeded"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// Cleanup removes old limiters to prevent memory leaks
func (rl *RateLimiter) Cleanup() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			// Remove limiters that haven't been used recently
			for clientID, limiter := range rl.limiters {
				if limiter.Tokens() == float64(rl.burst) {
					delete(rl.limiters, clientID)
				}
			}
			rl.mu.Unlock()
		}
	}()
}
