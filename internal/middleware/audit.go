package middleware

import (
	"bytes"
	"io"
	"time"

	"healthcare-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// AuditMiddleware logs all API requests for compliance
type AuditMiddleware struct {
	repo   *repository.BaseRepository
	logger *logrus.Logger
}

// NewAuditMiddleware creates a new audit middleware
func NewAuditMiddleware(repo *repository.BaseRepository, logger *logrus.Logger) *AuditMiddleware {
	return &AuditMiddleware{
		repo:   repo,
		logger: logger,
	}
}

// AuditLog middleware logs all requests for healthcare compliance
func (am *AuditMiddleware) AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Capture request body for audit
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Get user info from context (if authenticated)
		userID, _ := c.Get("user_id")
		userIDStr, _ := userID.(string)

		// Process request
		c.Next()

		// Log audit entry
		duration := time.Since(start)
		
		auditEntry := map[string]interface{}{
			"request_id":    requestID,
			"timestamp":     start.UTC(),
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"query":         c.Request.URL.RawQuery,
			"status_code":   c.Writer.Status(),
			"duration_ms":   duration.Milliseconds(),
			"client_ip":     c.ClientIP(),
			"user_agent":    c.Request.UserAgent(),
			"user_id":       userIDStr,
			"request_size":  len(requestBody),
			"response_size": c.Writer.Size(),
		}

		// Log sensitive operations with more detail
		if c.Request.Method != "GET" {
			auditEntry["request_body"] = string(requestBody)
		}

		am.logger.WithFields(auditEntry).Info("API Request Audit")

		// Store in database for compliance (async)
		go func() {
			// Implementation would store audit log in database
			// This is important for healthcare compliance (HIPAA, etc.)
		}()
	}
}
