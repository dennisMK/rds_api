package service

import (
	"context"
	"fmt"
	"time"

	"healthcare-api/internal/models"
	"healthcare-api/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ObservationService struct {
	repo   *repository.ObservationRepository
	logger *logrus.Logger
}

func NewObservationService(repo *repository.ObservationRepository, logger *logrus.Logger) *ObservationService {
	return &ObservationService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ObservationService) CreateObservation(ctx context.Context, req *models.ObservationCreateRequest) (*models.Observation, error) {
	s.logger.WithContext(ctx).Info("Creating new observation")

	// Generate UUID for new observation
	observationID := uuid.New()

	// Convert request to observation model
	observation := &models.Observation{
		Resource: models.Resource{
			ID:        observationID,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Version:   1,
		},
		Identifier:           req.Identifier,
		BasedOn:              req.BasedOn,
		PartOf:               req.PartOf,
		Status:               req.Status,
		Category:             req.Category,
		Code:                 req.Code,
		Subject:              req.Subject,
		Focus:                req.Focus,
		Encounter:            req.Encounter,
		EffectiveDateTime:    req.EffectiveDateTime,
		EffectivePeriod:      req.EffectivePeriod,
		EffectiveTiming:      req.EffectiveTiming,
		EffectiveInstant:     req.EffectiveInstant,
		Issued:               req.Issued,
		Performer:            req.Performer,
		ValueQuantity:        req.ValueQuantity,
		ValueCodeableConcept: req.ValueCodeableConcept,
		ValueString:          req.ValueString,
		ValueBoolean:         req.ValueBoolean,
		ValueInteger:         req.ValueInteger,
		ValueRange:           req.ValueRange,
		ValueRatio:           req.ValueRatio,
		ValueSampledData:     req.ValueSampledData,
		ValueTime:            req.ValueTime,
		ValueDateTime:        req.ValueDateTime,
		ValuePeriod:          req.ValuePeriod,
		DataAbsentReason:     req.DataAbsentReason,
		Interpretation:       req.Interpretation,
		Note:                 req.Note,
		BodySite:             req.BodySite,
		Method:               req.Method,
		Specimen:             req.Specimen,
		Device:               req.Device,
		ReferenceRange:       req.ReferenceRange,
		HasMember:            req.HasMember,
		DerivedFrom:          req.DerivedFrom,
		Component:            req.Component,
	}

	// Create observation in repository
	if err := s.repo.Create(ctx, observation); err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("Failed to create observation")
		return nil, fmt.Errorf("failed to create observation: %w", err)
	}

	s.logger.WithContext(ctx).WithField("observation_id", observation.ID).Info("Observation created successfully")
	return observation, nil
}

func (s *ObservationService) GetObservation(ctx context.Context, id uuid.UUID) (*models.Observation, error) {
	s.logger.WithContext(ctx).WithField("observation_id", id).Info("Retrieving observation")

	observation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).WithField("observation_id", id).Error("Failed to retrieve observation")
		return nil, fmt.Errorf("failed to retrieve observation: %w", err)
	}

	return observation, nil
}

func (s *ObservationService) UpdateObservation(ctx context.Context, id uuid.UUID, req *models.ObservationUpdateRequest) (*models.Observation, error) {
	s.logger.WithContext(ctx).WithField("observation_id", id).Info("Updating observation")

	// Get existing observation
	existingObservation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing observation: %w", err)
	}

	// Update fields that are provided in the request
	if req.Identifier != nil {
		existingObservation.Identifier = req.Identifier
	}
	if req.BasedOn != nil {
		existingObservation.BasedOn = req.BasedOn
	}
	if req.PartOf != nil {
		existingObservation.PartOf = req.PartOf
	}
	if req.Status != nil {
		existingObservation.Status = *req.Status
	}
	if req.Category != nil {
		existingObservation.Category = req.Category
	}
	if req.Code != nil {
		existingObservation.Code = *req.Code
	}
	if req.Subject != nil {
		existingObservation.Subject = *req.Subject
	}
	if req.Focus != nil {
		existingObservation.Focus = req.Focus
	}
	if req.Encounter != nil {
		existingObservation.Encounter = req.Encounter
	}
	if req.EffectiveDateTime != nil {
		existingObservation.EffectiveDateTime = req.EffectiveDateTime
	}
	if req.EffectivePeriod != nil {
		existingObservation.EffectivePeriod = req.EffectivePeriod
	}
	if req.EffectiveTiming != nil {
		existingObservation.EffectiveTiming = req.EffectiveTiming
	}
	if req.EffectiveInstant != nil {
		existingObservation.EffectiveInstant = req.EffectiveInstant
	}
	if req.Issued != nil {
		existingObservation.Issued = req.Issued
	}
	if req.Performer != nil {
		existingObservation.Performer = req.Performer
	}
	if req.ValueQuantity != nil {
		existingObservation.ValueQuantity = req.ValueQuantity
	}
	if req.ValueCodeableConcept != nil {
		existingObservation.ValueCodeableConcept = req.ValueCodeableConcept
	}
	if req.ValueString != nil {
		existingObservation.ValueString = req.ValueString
	}
	if req.ValueBoolean != nil {
		existingObservation.ValueBoolean = req.ValueBoolean
	}
	if req.ValueInteger != nil {
		existingObservation.ValueInteger = req.ValueInteger
	}
	if req.ValueRange != nil {
		existingObservation.ValueRange = req.ValueRange
	}
	if req.ValueRatio != nil {
		existingObservation.ValueRatio = req.ValueRatio
	}
	if req.ValueSampledData != nil {
		existingObservation.ValueSampledData = req.ValueSampledData
	}
	if req.ValueTime != nil {
		existingObservation.ValueTime = req.ValueTime
	}
	if req.ValueDateTime != nil {
		existingObservation.ValueDateTime = req.ValueDateTime
	}
	if req.ValuePeriod != nil {
		existingObservation.ValuePeriod = req.ValuePeriod
	}
	if req.DataAbsentReason != nil {
		existingObservation.DataAbsentReason = req.DataAbsentReason
	}
	if req.Interpretation != nil {
		existingObservation.Interpretation = req.Interpretation
	}
	if req.Note != nil {
		existingObservation.Note = req.Note
	}
	if req.BodySite != nil {
		existingObservation.BodySite = req.BodySite
	}
	if req.Method != nil {
		existingObservation.Method = req.Method
	}
	if req.Specimen != nil {
		existingObservation.Specimen = req.Specimen
	}
	if req.Device != nil {
		existingObservation.Device = req.Device
	}
	if req.ReferenceRange != nil {
		existingObservation.ReferenceRange = req.ReferenceRange
	}
	if req.HasMember != nil {
		existingObservation.HasMember = req.HasMember
	}
	if req.DerivedFrom != nil {
		existingObservation.DerivedFrom = req.DerivedFrom
	}
	if req.Component != nil {
		existingObservation.Component = req.Component
	}

	// Update in repository
	if err := s.repo.Update(ctx, existingObservation); err != nil {
		s.logger.WithContext(ctx).WithError(err).WithField("observation_id", id).Error("Failed to update observation")
		return nil, fmt.Errorf("failed to update observation: %w", err)
	}

	s.logger.WithContext(ctx).WithField("observation_id", id).Info("Observation updated successfully")
	return existingObservation, nil
}

func (s *ObservationService) DeleteObservation(ctx context.Context, id uuid.UUID) error {
	s.logger.WithContext(ctx).WithField("observation_id", id).Info("Deleting observation")

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.WithContext(ctx).WithError(err).WithField("observation_id", id).Error("Failed to delete observation")
		return fmt.Errorf("failed to delete observation: %w", err)
	}

	s.logger.WithContext(ctx).WithField("observation_id", id).Info("Observation deleted successfully")
	return nil
}

func (s *ObservationService) ListObservations(ctx context.Context, limit, offset int) (*models.ObservationListResponse, error) {
	s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"limit":  limit,
		"offset": offset,
	}).Info("Listing observations")

	// Validate and set pagination parameters
	params := repository.ValidatePaginationParams(limit, offset)

	observations, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("Failed to list observations")
		return nil, fmt.Errorf("failed to list observations: %w", err)
	}

	// Convert to response format
	entries := make([]models.ObservationEntry, len(observations))
	for i, observation := range observations {
		entries[i] = models.ObservationEntry{
			FullURL:  fmt.Sprintf("/api/v1/observations/%s", observation.ID),
			Resource: observation,
			Search: &models.SearchEntry{
				Mode: "match",
			},
		}
	}

	response := &models.ObservationListResponse{
		ResourceType: "Bundle",
		ID:           uuid.New().String(),
		Type:         "searchset",
		Total:        pagination.Total,
		Entry:        entries,
	}

	// Add pagination links
	if pagination.HasNext {
		response.Link = append(response.Link, models.BundleLink{
			Relation: "next",
			URL:      fmt.Sprintf("/api/v1/observations?limit=%d&offset=%d", params.Limit, params.Offset+params.Limit),
		})
	}

	if params.Offset > 0 {
		prevOffset := params.Offset - params.Limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		response.Link = append(response.Link, models.BundleLink{
			Relation: "prev",
			URL:      fmt.Sprintf("/api/v1/observations?limit=%d&offset=%d", params.Limit, prevOffset),
		})
	}

	s.logger.WithContext(ctx).WithField("total", pagination.Total).Info("Observations listed successfully")
	return response, nil
}
