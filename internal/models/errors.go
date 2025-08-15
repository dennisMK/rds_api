package models

import (
	"fmt"
	"net/http"
)

// APIError represents a standardized API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.Code, e.Message)
}

// Common API errors
var (
	ErrBadRequest = APIError{
		Code:    http.StatusBadRequest,
		Message: "Bad Request",
	}
	
	ErrUnauthorized = APIError{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
	}
	
	ErrForbidden = APIError{
		Code:    http.StatusForbidden,
		Message: "Forbidden",
	}
	
	ErrNotFound = APIError{
		Code:    http.StatusNotFound,
		Message: "Resource Not Found",
	}
	
	ErrConflict = APIError{
		Code:    http.StatusConflict,
		Message: "Resource Conflict",
	}
	
	ErrUnprocessableEntity = APIError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Unprocessable Entity",
	}
	
	ErrInternalServer = APIError{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
	}
	
	ErrServiceUnavailable = APIError{
		Code:    http.StatusServiceUnavailable,
		Message: "Service Unavailable",
	}
)

// NewAPIError creates a new API error with custom details
func NewAPIError(code int, message, details string) APIError {
	return APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (v ValidationErrors) Error() string {
	return fmt.Sprintf("Validation failed with %d errors", len(v.Errors))
}

// OperationOutcome represents a FHIR OperationOutcome resource
type OperationOutcome struct {
	ResourceType string                    `json:"resourceType"`
	ID           string                    `json:"id,omitempty"`
	Meta         *Meta                     `json:"meta,omitempty"`
	Issue        []OperationOutcomeIssue   `json:"issue"`
}

// OperationOutcomeIssue represents an issue in an OperationOutcome
type OperationOutcomeIssue struct {
	Severity    string           `json:"severity" validate:"required,oneof=fatal error warning information"`
	Code        string           `json:"code" validate:"required"`
	Details     *CodeableConcept `json:"details,omitempty"`
	Diagnostics *string          `json:"diagnostics,omitempty"`
	Location    []string         `json:"location,omitempty"`
	Expression  []string         `json:"expression,omitempty"`
}

// NewOperationOutcome creates a new OperationOutcome
func NewOperationOutcome(severity, code, diagnostics string) *OperationOutcome {
	return &OperationOutcome{
		ResourceType: "OperationOutcome",
		Issue: []OperationOutcomeIssue{
			{
				Severity:    severity,
				Code:        code,
				Diagnostics: &diagnostics,
			},
		},
	}
}
