package middleware

import (
	"net/http"

	"healthcare-api/internal/models"
	"healthcare-api/internal/validation"

	"github.com/gin-gonic/gin"
)

// ValidationMiddleware provides request validation
type ValidationMiddleware struct {
	validator *validation.Validator
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		validator: validation.NewValidator(),
	}
}

// ValidatePatientCreate validates patient creation requests
func (vm *ValidationMiddleware) ValidatePatientCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PatientCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid JSON: "+err.Error()))
			c.Abort()
			return
		}

		if validationErrors := vm.validator.ValidatePatientCreate(&req); validationErrors != nil {
			outcome := models.NewOperationOutcome("error", "invalid", "Validation failed")
			for _, validationError := range validationErrors.Errors {
				outcome.Issue = append(outcome.Issue, models.OperationOutcomeIssue{
					Severity:    "error",
					Code:        "invalid",
					Diagnostics: &validationError.Message,
					Expression:  []string{validationError.Field},
				})
			}
			c.JSON(http.StatusUnprocessableEntity, outcome)
			c.Abort()
			return
		}

		// Store validated request in context
		c.Set("validated_request", &req)
		c.Next()
	}
}

// ValidatePatientUpdate validates patient update requests
func (vm *ValidationMiddleware) ValidatePatientUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PatientUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid JSON: "+err.Error()))
			c.Abort()
			return
		}

		if validationErrors := vm.validator.ValidatePatientUpdate(&req); validationErrors != nil {
			outcome := models.NewOperationOutcome("error", "invalid", "Validation failed")
			for _, validationError := range validationErrors.Errors {
				outcome.Issue = append(outcome.Issue, models.OperationOutcomeIssue{
					Severity:    "error",
					Code:        "invalid",
					Diagnostics: &validationError.Message,
					Expression:  []string{validationError.Field},
				})
			}
			c.JSON(http.StatusUnprocessableEntity, outcome)
			c.Abort()
			return
		}

		c.Set("validated_request", &req)
		c.Next()
	}
}

// ValidateObservationCreate validates observation creation requests
func (vm *ValidationMiddleware) ValidateObservationCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ObservationCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid JSON: "+err.Error()))
			c.Abort()
			return
		}

		if validationErrors := vm.validator.ValidateObservationCreate(&req); validationErrors != nil {
			outcome := models.NewOperationOutcome("error", "invalid", "Validation failed")
			for _, validationError := range validationErrors.Errors {
				outcome.Issue = append(outcome.Issue, models.OperationOutcomeIssue{
					Severity:    "error",
					Code:        "invalid",
					Diagnostics: &validationError.Message,
					Expression:  []string{validationError.Field},
				})
			}
			c.JSON(http.StatusUnprocessableEntity, outcome)
			c.Abort()
			return
		}

		c.Set("validated_request", &req)
		c.Next()
	}
}

// ValidateObservationUpdate validates observation update requests
func (vm *ValidationMiddleware) ValidateObservationUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ObservationUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid JSON: "+err.Error()))
			c.Abort()
			return
		}

		if validationErrors := vm.validator.ValidateObservationUpdate(&req); validationErrors != nil {
			outcome := models.NewOperationOutcome("error", "invalid", "Validation failed")
			for _, validationError := range validationErrors.Errors {
				outcome.Issue = append(outcome.Issue, models.OperationOutcomeIssue{
					Severity:    "error",
					Code:        "invalid",
					Diagnostics: &validationError.Message,
					Expression:  []string{validationError.Field},
				})
			}
			c.JSON(http.StatusUnprocessableEntity, outcome)
			c.Abort()
			return
		}

		c.Set("validated_request", &req)
		c.Next()
	}
}
