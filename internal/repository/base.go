package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"healthcare-api/internal/database"

	"github.com/google/uuid"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	db *database.DB
}

func NewBaseRepository(db *database.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           uuid.UUID       `json:"id"`
	ResourceType string          `json:"resource_type"`
	ResourceID   uuid.UUID       `json:"resource_id"`
	Action       string          `json:"action"`
	UserID       *string         `json:"user_id,omitempty"`
	UserAgent    *string         `json:"user_agent,omitempty"`
	IPAddress    *string         `json:"ip_address,omitempty"`
	RequestID    *string         `json:"request_id,omitempty"`
	OldValues    json.RawMessage `json:"old_values,omitempty"`
	NewValues    json.RawMessage `json:"new_values,omitempty"`
	Timestamp    time.Time       `json:"timestamp"`
}

// LogAudit creates an audit log entry
func (r *BaseRepository) LogAudit(ctx context.Context, log *AuditLog) error {
	query := `
		INSERT INTO audit_logs (resource_type, resource_id, action, user_id, user_agent, ip_address, request_id, old_values, new_values)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ResourceType,
		log.ResourceID,
		log.Action,
		log.UserID,
		log.UserAgent,
		log.IPAddress,
		log.RequestID,
		log.OldValues,
		log.NewValues,
	)

	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// PaginationResult represents paginated results
type PaginationResult struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	HasNext bool `json:"has_next"`
}

// GetPaginationResult calculates pagination metadata
func GetPaginationResult(total int64, params PaginationParams) PaginationResult {
	hasNext := int64(params.Offset+params.Limit) < total
	
	return PaginationResult{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasNext: hasNext,
	}
}

// ValidatePaginationParams validates and sets default pagination parameters
func ValidatePaginationParams(limit, offset int) PaginationParams {
	if limit <= 0 || limit > 100 {
		limit = 20 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	return PaginationParams{
		Limit:  limit,
		Offset: offset,
	}
}
