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

type PatientService struct {
	repo   *repository.PatientRepository
	logger *logrus.Logger
}

func NewPatientService(repo *repository.PatientRepository, logger *logrus.Logger) *PatientService {
	return &PatientService{
		repo:   repo,
		logger: logger,
	}
}

func (s *PatientService) CreatePatient(ctx context.Context, req *models.PatientCreateRequest) (*models.Patient, error) {
	s.logger.WithContext(ctx).Info("Creating new patient")

	// Generate UUID for new patient
	patientID := uuid.New()

	// Convert request to patient model
	patient := &models.Patient{
		Resource: models.Resource{
			ID:        patientID,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Version:   1,
		},
		Identifier:              req.Identifier,
		Active:                  req.Active,
		Name:                    req.Name,
		Telecom:                 req.Telecom,
		Gender:                  req.Gender,
		BirthDate:               req.BirthDate,
		DeceasedBoolean:         req.DeceasedBoolean,
		DeceasedDateTime:        req.DeceasedDateTime,
		Address:                 req.Address,
		MaritalStatus:           req.MaritalStatus,
		MultipleBirthBoolean:    req.MultipleBirthBoolean,
		MultipleBirthInteger:    req.MultipleBirthInteger,
		Photo:                   req.Photo,
		Contact:                 req.Contact,
		Communication:           req.Communication,
		GeneralPractitioner:     req.GeneralPractitioner,
		ManagingOrganization:    req.ManagingOrganization,
		Link:                    req.Link,
	}

	// Set default active status if not provided
	if patient.Active == nil {
		active := true
		patient.Active = &active
	}

	// Create patient in repository
	if err := s.repo.Create(ctx, patient); err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("Failed to create patient")
		return nil, fmt.Errorf("failed to create patient: %w", err)
	}

	s.logger.WithContext(ctx).WithField("patient_id", patient.ID).Info("Patient created successfully")
	return patient, nil
}

func (s *PatientService) GetPatient(ctx context.Context, id uuid.UUID) (*models.Patient, error) {
	s.logger.WithContext(ctx).WithField("patient_id", id).Info("Retrieving patient")

	patient, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).WithField("patient_id", id).Error("Failed to retrieve patient")
		return nil, fmt.Errorf("failed to retrieve patient: %w", err)
	}

	return patient, nil
}

func (s *PatientService) UpdatePatient(ctx context.Context, id uuid.UUID, req *models.PatientUpdateRequest) (*models.Patient, error) {
	s.logger.WithContext(ctx).WithField("patient_id", id).Info("Updating patient")

	// Get existing patient
	existingPatient, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing patient: %w", err)
	}

	// Update fields that are provided in the request
	if req.Identifier != nil {
		existingPatient.Identifier = req.Identifier
	}
	if req.Active != nil {
		existingPatient.Active = req.Active
	}
	if req.Name != nil {
		existingPatient.Name = req.Name
	}
	if req.Telecom != nil {
		existingPatient.Telecom = req.Telecom
	}
	if req.Gender != nil {
		existingPatient.Gender = req.Gender
	}
	if req.BirthDate != nil {
		existingPatient.BirthDate = req.BirthDate
	}
	if req.DeceasedBoolean != nil {
		existingPatient.DeceasedBoolean = req.DeceasedBoolean
	}
	if req.DeceasedDateTime != nil {
		existingPatient.DeceasedDateTime = req.DeceasedDateTime
	}
	if req.Address != nil {
		existingPatient.Address = req.Address
	}
	if req.MaritalStatus != nil {
		existingPatient.MaritalStatus = req.MaritalStatus
	}
	if req.MultipleBirthBoolean != nil {
		existingPatient.MultipleBirthBoolean = req.MultipleBirthBoolean
	}
	if req.MultipleBirthInteger != nil {
		existingPatient.MultipleBirthInteger = req.MultipleBirthInteger
	}
	if req.Photo != nil {
		existingPatient.Photo = req.Photo
	}
	if req.Contact != nil {
		existingPatient.Contact = req.Contact
	}
	if req.Communication != nil {
		existingPatient.Communication = req.Communication
	}
	if req.GeneralPractitioner != nil {
		existingPatient.GeneralPractitioner = req.GeneralPractitioner
	}
	if req.ManagingOrganization != nil {
		existingPatient.ManagingOrganization = req.ManagingOrganization
	}
	if req.Link != nil {
		existingPatient.Link = req.Link
	}

	// Update in repository
	if err := s.repo.Update(ctx, existingPatient); err != nil {
		s.logger.WithContext(ctx).WithError(err).WithField("patient_id", id).Error("Failed to update patient")
		return nil, fmt.Errorf("failed to update patient: %w", err)
	}

	s.logger.WithContext(ctx).WithField("patient_id", id).Info("Patient updated successfully")
	return existingPatient, nil
}

func (s *PatientService) DeletePatient(ctx context.Context, id uuid.UUID) error {
	s.logger.WithContext(ctx).WithField("patient_id", id).Info("Deleting patient")

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.WithContext(ctx).WithError(err).WithField("patient_id", id).Error("Failed to delete patient")
		return fmt.Errorf("failed to delete patient: %w", err)
	}

	s.logger.WithContext(ctx).WithField("patient_id", id).Info("Patient deleted successfully")
	return nil
}

func (s *PatientService) ListPatients(ctx context.Context, limit, offset int) (*models.PatientListResponse, error) {
	s.logger.WithContext(ctx).WithFields(logrus.Fields{
		"limit":  limit,
		"offset": offset,
	}).Info("Listing patients")

	// Validate and set pagination parameters
	params := repository.ValidatePaginationParams(limit, offset)

	patients, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("Failed to list patients")
		return nil, fmt.Errorf("failed to list patients: %w", err)
	}

	// Convert to response format
	entries := make([]models.PatientEntry, len(patients))
	for i, patient := range patients {
		entries[i] = models.PatientEntry{
			FullURL:  fmt.Sprintf("/api/v1/patients/%s", patient.ID),
			Resource: patient,
			Search: &models.SearchEntry{
				Mode: "match",
			},
		}
	}

	response := &models.PatientListResponse{
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
			URL:      fmt.Sprintf("/api/v1/patients?limit=%d&offset=%d", params.Limit, params.Offset+params.Limit),
		})
	}

	if params.Offset > 0 {
		prevOffset := params.Offset - params.Limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		response.Link = append(response.Link, models.BundleLink{
			Relation: "prev",
			URL:      fmt.Sprintf("/api/v1/patients?limit=%d&offset=%d", params.Limit, prevOffset),
		})
	}

	s.logger.WithContext(ctx).WithField("total", pagination.Total).Info("Patients listed successfully")
	return response, nil
}
