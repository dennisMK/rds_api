package handlers

import (
	"net/http"
	"strconv"

	"healthcare-api/internal/models"
	"healthcare-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PatientHandler struct {
	service *service.PatientService
	logger  *logrus.Logger
}

func NewPatientHandler(service *service.PatientService, logger *logrus.Logger) *PatientHandler {
	return &PatientHandler{
		service: service,
		logger:  logger,
	}
}

// CreatePatient handles POST /api/v1/patients
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var req models.PatientCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind patient create request")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid request body: "+err.Error()))
		return
	}

	patient, err := h.service.CreatePatient(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create patient")
		c.JSON(http.StatusInternalServerError, models.NewOperationOutcome("error", "exception", "Failed to create patient"))
		return
	}

	c.Header("Location", "/api/v1/patients/"+patient.ID.String())
	c.JSON(http.StatusCreated, patient)
}

// GetPatient handles GET /api/v1/patients/:id
func (h *PatientHandler) GetPatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).WithField("id", idStr).Error("Invalid patient ID")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid patient ID format"))
		return
	}

	patient, err := h.service.GetPatient(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to get patient")
		if err.Error() == "patient not found" {
			c.JSON(http.StatusNotFound, models.NewOperationOutcome("error", "not-found", "Patient not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewOperationOutcome("error", "exception", "Failed to retrieve patient"))
		return
	}

	c.JSON(http.StatusOK, patient)
}

// UpdatePatient handles PUT /api/v1/patients/:id
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).WithField("id", idStr).Error("Invalid patient ID")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid patient ID format"))
		return
	}

	var req models.PatientUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind patient update request")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid request body: "+err.Error()))
		return
	}

	patient, err := h.service.UpdatePatient(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to update patient")
		if err.Error() == "patient not found" {
			c.JSON(http.StatusNotFound, models.NewOperationOutcome("error", "not-found", "Patient not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewOperationOutcome("error", "exception", "Failed to update patient"))
		return
	}

	c.JSON(http.StatusOK, patient)
}

// DeletePatient handles DELETE /api/v1/patients/:id
func (h *PatientHandler) DeletePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).WithField("id", idStr).Error("Invalid patient ID")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid patient ID format"))
		return
	}

	err = h.service.DeletePatient(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to delete patient")
		if err.Error() == "patient not found" {
			c.JSON(http.StatusNotFound, models.NewOperationOutcome("error", "not-found", "Patient not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewOperationOutcome("error", "exception", "Failed to delete patient"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListPatients handles GET /api/v1/patients
func (h *PatientHandler) ListPatients(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.logger.WithError(err).WithField("limit", limitStr).Error("Invalid limit parameter")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid limit parameter"))
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.logger.WithError(err).WithField("offset", offsetStr).Error("Invalid offset parameter")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid offset parameter"))
		return
	}

	response, err := h.service.ListPatients(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list patients")
		c.JSON(http.StatusInternalServerError, models.NewOperationOutcome("error", "exception", "Failed to list patients"))
		return
	}

	c.JSON(http.StatusOK, response)
}
