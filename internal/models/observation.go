package models

import (
	"time"
)

// Observation represents a FHIR Observation resource
type Observation struct {
	Resource
	
	// Observation-specific fields
	Identifier           []Identifier      `json:"identifier,omitempty" db:"identifier"`
	BasedOn              []Reference       `json:"basedOn,omitempty" db:"based_on"`
	PartOf               []Reference       `json:"partOf,omitempty" db:"part_of"`
	Status               string            `json:"status" db:"status" validate:"required,oneof=registered preliminary final amended corrected cancelled entered-in-error unknown"`
	Category             []CodeableConcept `json:"category,omitempty" db:"category"`
	Code                 CodeableConcept   `json:"code" db:"code" validate:"required"`
	Subject              Reference         `json:"subject" db:"subject" validate:"required"`
	Focus                []Reference       `json:"focus,omitempty" db:"focus"`
	Encounter            *Reference        `json:"encounter,omitempty" db:"encounter"`
	EffectiveDateTime    *time.Time        `json:"effectiveDateTime,omitempty" db:"effective_date_time"`
	EffectivePeriod      *Period           `json:"effectivePeriod,omitempty" db:"effective_period"`
	EffectiveTiming      *Timing           `json:"effectiveTiming,omitempty" db:"effective_timing"`
	EffectiveInstant     *time.Time        `json:"effectiveInstant,omitempty" db:"effective_instant"`
	Issued               *time.Time        `json:"issued,omitempty" db:"issued"`
	Performer            []Reference       `json:"performer,omitempty" db:"performer"`
	ValueQuantity        *Quantity         `json:"valueQuantity,omitempty" db:"value_quantity"`
	ValueCodeableConcept *CodeableConcept  `json:"valueCodeableConcept,omitempty" db:"value_codeable_concept"`
	ValueString          *string           `json:"valueString,omitempty" db:"value_string"`
	ValueBoolean         *bool             `json:"valueBoolean,omitempty" db:"value_boolean"`
	ValueInteger         *int              `json:"valueInteger,omitempty" db:"value_integer"`
	ValueRange           *Range            `json:"valueRange,omitempty" db:"value_range"`
	ValueRatio           *Ratio            `json:"valueRatio,omitempty" db:"value_ratio"`
	ValueSampledData     *SampledData      `json:"valueSampledData,omitempty" db:"value_sampled_data"`
	ValueTime            *string           `json:"valueTime,omitempty" db:"value_time"`
	ValueDateTime        *time.Time        `json:"valueDateTime,omitempty" db:"value_date_time"`
	ValuePeriod          *Period           `json:"valuePeriod,omitempty" db:"value_period"`
	DataAbsentReason     *CodeableConcept  `json:"dataAbsentReason,omitempty" db:"data_absent_reason"`
	Interpretation       []CodeableConcept `json:"interpretation,omitempty" db:"interpretation"`
	Note                 []Annotation      `json:"note,omitempty" db:"note"`
	BodySite             *CodeableConcept  `json:"bodySite,omitempty" db:"body_site"`
	Method               *CodeableConcept  `json:"method,omitempty" db:"method"`
	Specimen             *Reference        `json:"specimen,omitempty" db:"specimen"`
	Device               *Reference        `json:"device,omitempty" db:"device"`
	ReferenceRange       []ObservationReferenceRange `json:"referenceRange,omitempty" db:"reference_range"`
	HasMember            []Reference       `json:"hasMember,omitempty" db:"has_member"`
	DerivedFrom          []Reference       `json:"derivedFrom,omitempty" db:"derived_from"`
	Component            []ObservationComponent `json:"component,omitempty" db:"component"`
}

// ObservationReferenceRange represents reference ranges for observations
type ObservationReferenceRange struct {
	Low           *Quantity        `json:"low,omitempty"`
	High          *Quantity        `json:"high,omitempty"`
	Type          *CodeableConcept `json:"type,omitempty"`
	AppliesTo     []CodeableConcept `json:"appliesTo,omitempty"`
	Age           *Range           `json:"age,omitempty"`
	Text          *string          `json:"text,omitempty"`
}

// ObservationComponent represents observation components
type ObservationComponent struct {
	Code                 CodeableConcept   `json:"code" validate:"required"`
	ValueQuantity        *Quantity         `json:"valueQuantity,omitempty"`
	ValueCodeableConcept *CodeableConcept  `json:"valueCodeableConcept,omitempty"`
	ValueString          *string           `json:"valueString,omitempty"`
	ValueBoolean         *bool             `json:"valueBoolean,omitempty"`
	ValueInteger         *int              `json:"valueInteger,omitempty"`
	ValueRange           *Range            `json:"valueRange,omitempty"`
	ValueRatio           *Ratio            `json:"valueRatio,omitempty"`
	ValueSampledData     *SampledData      `json:"valueSampledData,omitempty"`
	ValueTime            *string           `json:"valueTime,omitempty"`
	ValueDateTime        *time.Time        `json:"valueDateTime,omitempty"`
	ValuePeriod          *Period           `json:"valuePeriod,omitempty"`
	DataAbsentReason     *CodeableConcept  `json:"dataAbsentReason,omitempty"`
	Interpretation       []CodeableConcept `json:"interpretation,omitempty"`
	ReferenceRange       []ObservationReferenceRange `json:"referenceRange,omitempty"`
}

// Timing represents timing information
type Timing struct {
	Event  []time.Time   `json:"event,omitempty"`
	Repeat *TimingRepeat `json:"repeat,omitempty"`
	Code   *CodeableConcept `json:"code,omitempty"`
}

// TimingRepeat represents timing repeat information
type TimingRepeat struct {
	BoundsDuration *Duration `json:"boundsDuration,omitempty"`
	BoundsRange    *Range    `json:"boundsRange,omitempty"`
	BoundsPeriod   *Period   `json:"boundsPeriod,omitempty"`
	Count          *int      `json:"count,omitempty"`
	CountMax       *int      `json:"countMax,omitempty"`
	Duration       *float64  `json:"duration,omitempty"`
	DurationMax    *float64  `json:"durationMax,omitempty"`
	DurationUnit   *string   `json:"durationUnit,omitempty" validate:"omitempty,oneof=s min h d wk mo a"`
	Frequency      *int      `json:"frequency,omitempty"`
	FrequencyMax   *int      `json:"frequencyMax,omitempty"`
	Period         *float64  `json:"period,omitempty"`
	PeriodMax      *float64  `json:"periodMax,omitempty"`
	PeriodUnit     *string   `json:"periodUnit,omitempty" validate:"omitempty,oneof=s min h d wk mo a"`
	DayOfWeek      []string  `json:"dayOfWeek,omitempty"`
	TimeOfDay      []string  `json:"timeOfDay,omitempty"`
	When           []string  `json:"when,omitempty"`
	Offset         *int      `json:"offset,omitempty"`
}

// Duration represents a duration
type Duration struct {
	Value      *float64 `json:"value,omitempty"`
	Comparator *string  `json:"comparator,omitempty" validate:"omitempty,oneof=< <= >= > ad"`
	Unit       *string  `json:"unit,omitempty"`
	System     *string  `json:"system,omitempty" validate:"omitempty,uri"`
	Code       *string  `json:"code,omitempty"`
}

// ObservationCreateRequest represents the request to create an observation
type ObservationCreateRequest struct {
	Identifier           []Identifier      `json:"identifier,omitempty"`
	BasedOn              []Reference       `json:"basedOn,omitempty"`
	PartOf               []Reference       `json:"partOf,omitempty"`
	Status               string            `json:"status" validate:"required,oneof=registered preliminary final amended corrected cancelled entered-in-error unknown"`
	Category             []CodeableConcept `json:"category,omitempty"`
	Code                 CodeableConcept   `json:"code" validate:"required"`
	Subject              Reference         `json:"subject" validate:"required"`
	Focus                []Reference       `json:"focus,omitempty"`
	Encounter            *Reference        `json:"encounter,omitempty"`
	EffectiveDateTime    *time.Time        `json:"effectiveDateTime,omitempty"`
	EffectivePeriod      *Period           `json:"effectivePeriod,omitempty"`
	EffectiveTiming      *Timing           `json:"effectiveTiming,omitempty"`
	EffectiveInstant     *time.Time        `json:"effectiveInstant,omitempty"`
	Issued               *time.Time        `json:"issued,omitempty"`
	Performer            []Reference       `json:"performer,omitempty"`
	ValueQuantity        *Quantity         `json:"valueQuantity,omitempty"`
	ValueCodeableConcept *CodeableConcept  `json:"valueCodeableConcept,omitempty"`
	ValueString          *string           `json:"valueString,omitempty"`
	ValueBoolean         *bool             `json:"valueBoolean,omitempty"`
	ValueInteger         *int              `json:"valueInteger,omitempty"`
	ValueRange           *Range            `json:"valueRange,omitempty"`
	ValueRatio           *Ratio            `json:"valueRatio,omitempty"`
	ValueSampledData     *SampledData      `json:"valueSampledData,omitempty"`
	ValueTime            *string           `json:"valueTime,omitempty"`
	ValueDateTime        *time.Time        `json:"valueDateTime,omitempty"`
	ValuePeriod          *Period           `json:"valuePeriod,omitempty"`
	DataAbsentReason     *CodeableConcept  `json:"dataAbsentReason,omitempty"`
	Interpretation       []CodeableConcept `json:"interpretation,omitempty"`
	Note                 []Annotation      `json:"note,omitempty"`
	BodySite             *CodeableConcept  `json:"bodySite,omitempty"`
	Method               *CodeableConcept  `json:"method,omitempty"`
	Specimen             *Reference        `json:"specimen,omitempty"`
	Device               *Reference        `json:"device,omitempty"`
	ReferenceRange       []ObservationReferenceRange `json:"referenceRange,omitempty"`
	HasMember            []Reference       `json:"hasMember,omitempty"`
	DerivedFrom          []Reference       `json:"derivedFrom,omitempty"`
	Component            []ObservationComponent `json:"component,omitempty"`
}

// ObservationUpdateRequest represents the request to update an observation
type ObservationUpdateRequest struct {
	Identifier           []Identifier      `json:"identifier,omitempty"`
	BasedOn              []Reference       `json:"basedOn,omitempty"`
	PartOf               []Reference       `json:"partOf,omitempty"`
	Status               *string           `json:"status,omitempty" validate:"omitempty,oneof=registered preliminary final amended corrected cancelled entered-in-error unknown"`
	Category             []CodeableConcept `json:"category,omitempty"`
	Code                 *CodeableConcept  `json:"code,omitempty"`
	Subject              *Reference        `json:"subject,omitempty"`
	Focus                []Reference       `json:"focus,omitempty"`
	Encounter            *Reference        `json:"encounter,omitempty"`
	EffectiveDateTime    *time.Time        `json:"effectiveDateTime,omitempty"`
	EffectivePeriod      *Period           `json:"effectivePeriod,omitempty"`
	EffectiveTiming      *Timing           `json:"effectiveTiming,omitempty"`
	EffectiveInstant     *time.Time        `json:"effectiveInstant,omitempty"`
	Issued               *time.Time        `json:"issued,omitempty"`
	Performer            []Reference       `json:"performer,omitempty"`
	ValueQuantity        *Quantity         `json:"valueQuantity,omitempty"`
	ValueCodeableConcept *CodeableConcept  `json:"valueCodeableConcept,omitempty"`
	ValueString          *string           `json:"valueString,omitempty"`
	ValueBoolean         *bool             `json:"valueBoolean,omitempty"`
	ValueInteger         *int              `json:"valueInteger,omitempty"`
	ValueRange           *Range            `json:"valueRange,omitempty"`
	ValueRatio           *Ratio            `json:"valueRatio,omitempty"`
	ValueSampledData     *SampledData      `json:"valueSampledData,omitempty"`
	ValueTime            *string           `json:"valueTime,omitempty"`
	ValueDateTime        *time.Time        `json:"valueDateTime,omitempty"`
	ValuePeriod          *Period           `json:"valuePeriod,omitempty"`
	DataAbsentReason     *CodeableConcept  `json:"dataAbsentReason,omitempty"`
	Interpretation       []CodeableConcept `json:"interpretation,omitempty"`
	Note                 []Annotation      `json:"note,omitempty"`
	BodySite             *CodeableConcept  `json:"bodySite,omitempty"`
	Method               *CodeableConcept  `json:"method,omitempty"`
	Specimen             *Reference        `json:"specimen,omitempty"`
	Device               *Reference        `json:"device,omitempty"`
	ReferenceRange       []ObservationReferenceRange `json:"referenceRange,omitempty"`
	HasMember            []Reference       `json:"hasMember,omitempty"`
	DerivedFrom          []Reference       `json:"derivedFrom,omitempty"`
	Component            []ObservationComponent `json:"component,omitempty"`
}

// ObservationListResponse represents the response for listing observations
type ObservationListResponse struct {
	ResourceType string           `json:"resourceType"`
	ID           string           `json:"id"`
	Type         string           `json:"type"`
	Total        int64            `json:"total"`
	Entry        []ObservationEntry `json:"entry"`
	Link         []BundleLink     `json:"link,omitempty"`
}

// ObservationEntry represents an observation entry in a bundle
type ObservationEntry struct {
	FullURL  string       `json:"fullUrl"`
	Resource *Observation `json:"resource"`
	Search   *SearchEntry `json:"search,omitempty"`
}
