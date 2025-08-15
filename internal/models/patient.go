package models

import (
	"time"
)

// Patient represents a FHIR Patient resource
type Patient struct {
	Resource
	
	// Patient-specific fields
	Identifier              []Identifier      `json:"identifier,omitempty" db:"identifier"`
	Active                  *bool             `json:"active,omitempty" db:"active"`
	Name                    []HumanName       `json:"name,omitempty" db:"name" validate:"required,min=1"`
	Telecom                 []ContactPoint    `json:"telecom,omitempty" db:"telecom"`
	Gender                  *string           `json:"gender,omitempty" db:"gender" validate:"omitempty,oneof=male female other unknown"`
	BirthDate               *time.Time        `json:"birthDate,omitempty" db:"birth_date"`
	DeceasedBoolean         *bool             `json:"deceasedBoolean,omitempty" db:"deceased_boolean"`
	DeceasedDateTime        *time.Time        `json:"deceasedDateTime,omitempty" db:"deceased_date_time"`
	Address                 []Address         `json:"address,omitempty" db:"address"`
	MaritalStatus           *CodeableConcept  `json:"maritalStatus,omitempty" db:"marital_status"`
	MultipleBirthBoolean    *bool             `json:"multipleBirthBoolean,omitempty" db:"multiple_birth_boolean"`
	MultipleBirthInteger    *int              `json:"multipleBirthInteger,omitempty" db:"multiple_birth_integer"`
	Photo                   []Attachment      `json:"photo,omitempty" db:"photo"`
	Contact                 []PatientContact  `json:"contact,omitempty" db:"contact"`
	Communication           []PatientCommunication `json:"communication,omitempty" db:"communication"`
	GeneralPractitioner     []Reference       `json:"generalPractitioner,omitempty" db:"general_practitioner"`
	ManagingOrganization    *Reference        `json:"managingOrganization,omitempty" db:"managing_organization"`
	Link                    []PatientLink     `json:"link,omitempty" db:"link"`
}

// PatientContact represents patient contact information
type PatientContact struct {
	Relationship    []CodeableConcept `json:"relationship,omitempty"`
	Name            *HumanName        `json:"name,omitempty"`
	Telecom         []ContactPoint    `json:"telecom,omitempty"`
	Address         *Address          `json:"address,omitempty"`
	Gender          *string           `json:"gender,omitempty" validate:"omitempty,oneof=male female other unknown"`
	Organization    *Reference        `json:"organization,omitempty"`
	Period          *Period           `json:"period,omitempty"`
}

// PatientCommunication represents patient communication preferences
type PatientCommunication struct {
	Language  CodeableConcept `json:"language" validate:"required"`
	Preferred *bool           `json:"preferred,omitempty"`
}

// PatientLink represents links to other patient resources
type PatientLink struct {
	Other Reference `json:"other" validate:"required"`
	Type  string    `json:"type" validate:"required,oneof=replaced-by replaces refer seealso"`
}

// PatientCreateRequest represents the request to create a patient
type PatientCreateRequest struct {
	Identifier              []Identifier      `json:"identifier,omitempty"`
	Active                  *bool             `json:"active,omitempty"`
	Name                    []HumanName       `json:"name" validate:"required,min=1"`
	Telecom                 []ContactPoint    `json:"telecom,omitempty"`
	Gender                  *string           `json:"gender,omitempty" validate:"omitempty,oneof=male female other unknown"`
	BirthDate               *time.Time        `json:"birthDate,omitempty"`
	DeceasedBoolean         *bool             `json:"deceasedBoolean,omitempty"`
	DeceasedDateTime        *time.Time        `json:"deceasedDateTime,omitempty"`
	Address                 []Address         `json:"address,omitempty"`
	MaritalStatus           *CodeableConcept  `json:"maritalStatus,omitempty"`
	MultipleBirthBoolean    *bool             `json:"multipleBirthBoolean,omitempty"`
	MultipleBirthInteger    *int              `json:"multipleBirthInteger,omitempty"`
	Photo                   []Attachment      `json:"photo,omitempty"`
	Contact                 []PatientContact  `json:"contact,omitempty"`
	Communication           []PatientCommunication `json:"communication,omitempty"`
	GeneralPractitioner     []Reference       `json:"generalPractitioner,omitempty"`
	ManagingOrganization    *Reference        `json:"managingOrganization,omitempty"`
	Link                    []PatientLink     `json:"link,omitempty"`
}

// PatientUpdateRequest represents the request to update a patient
type PatientUpdateRequest struct {
	Identifier              []Identifier      `json:"identifier,omitempty"`
	Active                  *bool             `json:"active,omitempty"`
	Name                    []HumanName       `json:"name,omitempty"`
	Telecom                 []ContactPoint    `json:"telecom,omitempty"`
	Gender                  *string           `json:"gender,omitempty" validate:"omitempty,oneof=male female other unknown"`
	BirthDate               *time.Time        `json:"birthDate,omitempty"`
	DeceasedBoolean         *bool             `json:"deceasedBoolean,omitempty"`
	DeceasedDateTime        *time.Time        `json:"deceasedDateTime,omitempty"`
	Address                 []Address         `json:"address,omitempty"`
	MaritalStatus           *CodeableConcept  `json:"maritalStatus,omitempty"`
	MultipleBirthBoolean    *bool             `json:"multipleBirthBoolean,omitempty"`
	MultipleBirthInteger    *int              `json:"multipleBirthInteger,omitempty"`
	Photo                   []Attachment      `json:"photo,omitempty"`
	Contact                 []PatientContact  `json:"contact,omitempty"`
	Communication           []PatientCommunication `json:"communication,omitempty"`
	GeneralPractitioner     []Reference       `json:"generalPractitioner,omitempty"`
	ManagingOrganization    *Reference        `json:"managingOrganization,omitempty"`
	Link                    []PatientLink     `json:"link,omitempty"`
}

// PatientListResponse represents the response for listing patients
type PatientListResponse struct {
	ResourceType string    `json:"resourceType"`
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Total        int64     `json:"total"`
	Entry        []PatientEntry `json:"entry"`
	Link         []BundleLink   `json:"link,omitempty"`
}

// PatientEntry represents a patient entry in a bundle
type PatientEntry struct {
	FullURL  string   `json:"fullUrl"`
	Resource *Patient `json:"resource"`
	Search   *SearchEntry `json:"search,omitempty"`
}

// SearchEntry represents search metadata
type SearchEntry struct {
	Mode  string   `json:"mode"`
	Score *float64 `json:"score,omitempty"`
}

// BundleLink represents a link in a bundle
type BundleLink struct {
	Relation string `json:"relation"`
	URL      string `json:"url"`
}
