package validation

import (
	"fmt"
	"reflect"
	"strings"

	"healthcare-api/internal/models"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	validate := validator.New()
	
	// Register custom validation functions
	validate.RegisterValidation("fhir_status", validateFHIRStatus)
	validate.RegisterValidation("fhir_gender", validateFHIRGender)
	validate.RegisterValidation("fhir_name_use", validateFHIRNameUse)
	validate.RegisterValidation("fhir_contact_system", validateFHIRContactSystem)
	validate.RegisterValidation("fhir_address_use", validateFHIRAddressUse)
	
	// Use JSON tag names in error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	
	return &Validator{validate: validate}
}

// ValidateStruct validates a struct and returns validation errors
func (v *Validator) ValidateStruct(s interface{}) *models.ValidationErrors {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors []models.ValidationError
	
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, validationErr := range validationErrs {
			validationErrors = append(validationErrors, models.ValidationError{
				Field:   validationErr.Field(),
				Message: getValidationMessage(validationErr),
				Value:   validationErr.Value(),
			})
		}
	}

	return &models.ValidationErrors{Errors: validationErrors}
}

// getValidationMessage returns a human-readable validation message
func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
	case "uri":
		return fmt.Sprintf("%s must be a valid URI", err.Field())
	case "fhir_status":
		return fmt.Sprintf("%s must be a valid FHIR status", err.Field())
	case "fhir_gender":
		return fmt.Sprintf("%s must be a valid FHIR gender", err.Field())
	case "fhir_name_use":
		return fmt.Sprintf("%s must be a valid FHIR name use", err.Field())
	case "fhir_contact_system":
		return fmt.Sprintf("%s must be a valid FHIR contact system", err.Field())
	case "fhir_address_use":
		return fmt.Sprintf("%s must be a valid FHIR address use", err.Field())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}

// Custom validation functions for FHIR-specific fields

func validateFHIRStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{"registered", "preliminary", "final", "amended", "corrected", "cancelled", "entered-in-error", "unknown"}
	
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func validateFHIRGender(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	validGenders := []string{"male", "female", "other", "unknown"}
	
	for _, validGender := range validGenders {
		if gender == validGender {
			return true
		}
	}
	return false
}

func validateFHIRNameUse(fl validator.FieldLevel) bool {
	use := fl.Field().String()
	validUses := []string{"usual", "official", "temp", "nickname", "anonymous", "old", "maiden"}
	
	for _, validUse := range validUses {
		if use == validUse {
			return true
		}
	}
	return false
}

func validateFHIRContactSystem(fl validator.FieldLevel) bool {
	system := fl.Field().String()
	validSystems := []string{"phone", "fax", "email", "pager", "url", "sms", "other"}
	
	for _, validSystem := range validSystems {
		if system == validSystem {
			return true
		}
	}
	return false
}

func validateFHIRAddressUse(fl validator.FieldLevel) bool {
	use := fl.Field().String()
	validUses := []string{"home", "work", "temp", "old", "billing"}
	
	for _, validUse := range validUses {
		if use == validUse {
			return true
		}
	}
	return false
}

// ValidatePatientCreate validates patient creation request
func (v *Validator) ValidatePatientCreate(req *models.PatientCreateRequest) *models.ValidationErrors {
	return v.ValidateStruct(req)
}

// ValidatePatientUpdate validates patient update request
func (v *Validator) ValidatePatientUpdate(req *models.PatientUpdateRequest) *models.ValidationErrors {
	return v.ValidateStruct(req)
}

// ValidateObservationCreate validates observation creation request
func (v *Validator) ValidateObservationCreate(req *models.ObservationCreateRequest) *models.ValidationErrors {
	return v.ValidateStruct(req)
}

// ValidateObservationUpdate validates observation update request
func (v *Validator) ValidateObservationUpdate(req *models.ObservationUpdateRequest) *models.ValidationErrors {
	return v.ValidateStruct(req)
}
