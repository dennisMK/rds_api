package models

import (
	"time"

	"github.com/google/uuid"
)

// Base FHIR resource elements that are common to all resources
type Resource struct {
	ID                uuid.UUID         `json:"id" db:"id"`
	Meta              *Meta             `json:"meta,omitempty" db:"meta"`
	ImplicitRules     *string           `json:"implicitRules,omitempty" db:"implicit_rules"`
	Language          *string           `json:"language,omitempty" db:"language"`
	Text              *Narrative        `json:"text,omitempty" db:"text"`
	Contained         []Resource        `json:"contained,omitempty" db:"contained"`
	Extension         []Extension       `json:"extension,omitempty" db:"extension"`
	ModifierExtension []Extension       `json:"modifierExtension,omitempty" db:"modifier_extension"`
	CreatedAt         time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time         `json:"updatedAt" db:"updated_at"`
	Version           int               `json:"version" db:"version"`
}

// Meta contains metadata about a resource
type Meta struct {
	VersionID   *string    `json:"versionId,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Source      *string    `json:"source,omitempty"`
	Profile     []string   `json:"profile,omitempty"`
	Security    []Coding   `json:"security,omitempty"`
	Tag         []Coding   `json:"tag,omitempty"`
}

// Narrative contains human-readable text
type Narrative struct {
	Status string `json:"status" validate:"required,oneof=generated extensions additional empty"`
	Div    string `json:"div" validate:"required"`
}

// Extension represents FHIR extensions
type Extension struct {
	URL                string      `json:"url" validate:"required,uri"`
	ValueString        *string     `json:"valueString,omitempty"`
	ValueInteger       *int        `json:"valueInteger,omitempty"`
	ValueBoolean       *bool       `json:"valueBoolean,omitempty"`
	ValueDateTime      *time.Time  `json:"valueDateTime,omitempty"`
	ValueCodeableConcept *CodeableConcept `json:"valueCodeableConcept,omitempty"`
	Extension          []Extension `json:"extension,omitempty"`
}

// Identifier represents a business identifier
type Identifier struct {
	Use      *string          `json:"use,omitempty" validate:"omitempty,oneof=usual official temp secondary old"`
	Type     *CodeableConcept `json:"type,omitempty"`
	System   *string          `json:"system,omitempty" validate:"omitempty,uri"`
	Value    *string          `json:"value,omitempty"`
	Period   *Period          `json:"period,omitempty"`
	Assigner *Reference       `json:"assigner,omitempty"`
}

// CodeableConcept represents a concept with coding
type CodeableConcept struct {
	Coding []Coding `json:"coding,omitempty"`
	Text   *string  `json:"text,omitempty"`
}

// Coding represents a code from a terminology system
type Coding struct {
	System       *string `json:"system,omitempty" validate:"omitempty,uri"`
	Version      *string `json:"version,omitempty"`
	Code         *string `json:"code,omitempty"`
	Display      *string `json:"display,omitempty"`
	UserSelected *bool   `json:"userSelected,omitempty"`
}

// Reference represents a reference to another resource
type Reference struct {
	Reference  *string     `json:"reference,omitempty"`
	Type       *string     `json:"type,omitempty"`
	Identifier *Identifier `json:"identifier,omitempty"`
	Display    *string     `json:"display,omitempty"`
}

// Period represents a time period
type Period struct {
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

// ContactPoint represents contact information
type ContactPoint struct {
	System *string `json:"system,omitempty" validate:"omitempty,oneof=phone fax email pager url sms other"`
	Value  *string `json:"value,omitempty"`
	Use    *string `json:"use,omitempty" validate:"omitempty,oneof=home work temp old mobile"`
	Rank   *int    `json:"rank,omitempty" validate:"omitempty,min=1"`
	Period *Period `json:"period,omitempty"`
}

// Address represents an address
type Address struct {
	Use        *string  `json:"use,omitempty" validate:"omitempty,oneof=home work temp old billing"`
	Type       *string  `json:"type,omitempty" validate:"omitempty,oneof=postal physical both"`
	Text       *string  `json:"text,omitempty"`
	Line       []string `json:"line,omitempty"`
	City       *string  `json:"city,omitempty"`
	District   *string  `json:"district,omitempty"`
	State      *string  `json:"state,omitempty"`
	PostalCode *string  `json:"postalCode,omitempty"`
	Country    *string  `json:"country,omitempty"`
	Period     *Period  `json:"period,omitempty"`
}

// HumanName represents a human name
type HumanName struct {
	Use    *string  `json:"use,omitempty" validate:"omitempty,oneof=usual official temp nickname anonymous old maiden"`
	Text   *string  `json:"text,omitempty"`
	Family *string  `json:"family,omitempty"`
	Given  []string `json:"given,omitempty"`
	Prefix []string `json:"prefix,omitempty"`
	Suffix []string `json:"suffix,omitempty"`
	Period *Period  `json:"period,omitempty"`
}

// Attachment represents an attachment
type Attachment struct {
	ContentType *string    `json:"contentType,omitempty"`
	Language    *string    `json:"language,omitempty"`
	Data        *string    `json:"data,omitempty"`
	URL         *string    `json:"url,omitempty" validate:"omitempty,uri"`
	Size        *int       `json:"size,omitempty"`
	Hash        *string    `json:"hash,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Creation    *time.Time `json:"creation,omitempty"`
}

// Quantity represents a measured amount
type Quantity struct {
	Value      *float64 `json:"value,omitempty"`
	Comparator *string  `json:"comparator,omitempty" validate:"omitempty,oneof=< <= >= > ad"`
	Unit       *string  `json:"unit,omitempty"`
	System     *string  `json:"system,omitempty" validate:"omitempty,uri"`
	Code       *string  `json:"code,omitempty"`
}

// Range represents a range of values
type Range struct {
	Low  *Quantity `json:"low,omitempty"`
	High *Quantity `json:"high,omitempty"`
}

// Ratio represents a ratio of two quantities
type Ratio struct {
	Numerator   *Quantity `json:"numerator,omitempty"`
	Denominator *Quantity `json:"denominator,omitempty"`
}

// SampledData represents sampled data
type SampledData struct {
	Origin     Quantity `json:"origin" validate:"required"`
	Period     float64  `json:"period" validate:"required"`
	Factor     *float64 `json:"factor,omitempty"`
	LowerLimit *float64 `json:"lowerLimit,omitempty"`
	UpperLimit *float64 `json:"upperLimit,omitempty"`
	Dimensions int      `json:"dimensions" validate:"required,min=1"`
	Data       *string  `json:"data,omitempty"`
}

// Annotation represents an annotation
type Annotation struct {
	AuthorReference *Reference `json:"authorReference,omitempty"`
	AuthorString    *string    `json:"authorString,omitempty"`
	Time            *time.Time `json:"time,omitempty"`
	Text            string     `json:"text" validate:"required"`
}
