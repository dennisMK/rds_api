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

type ObservationHandler struct {
	service *service.ObservationService
	logger  *logrus.Logger
}

func NewObservationHandler(service *service.ObservationService, logger *logrus.Logger) *ObservationHandler {
	return &ObservationHandler{
		service: service,
		logger:  logger,
	}
}

// CreateObservation handles POST /api/v1/observations
func (h *ObservationHandler) CreateObservation(c *gin.Context) {
	var req models.ObservationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind observation create request")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid request body: "+err.Error()))
		return
	}

	observation, err := h.service.CreateObservation(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create observation")
		c.JSON(http.StatusInternalServerError, models.NewOperationOutcome("error", "exception", "Failed to create observation"))
		return
	}

	c.Header("Location", "/api/v1/observations/"+observation.ID.String())
	c.JSON(http.StatusCreated, observation)
}

// GetObservation handles GET /api/v1/observations/:id
func (h *ObservationHandler) GetObservation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).WithField("id", idStr).Error("Invalid observation ID")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid observation ID format"))
		return
	}

	observation, err := h.service.GetObservation(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to get observation")
		if err.Error() == "observation not found" {
			c.JSON(http.StatusNotFound, models.NewOperationOutcome("error", "not-found", "Observation not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewOperationOutcome("error", "exception", "Failed to retrieve observation"))
		return
	}

	c.JSON(http.StatusOK, observation)
}

// UpdateObservation handles PUT /api/v1/observations/:id
func (h *ObservationHandler) UpdateObservation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).WithField("id", idStr).Error("Invalid observation ID")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid observation ID format"))
		return
	}

	var req models.ObservationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind observation update request")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid request body: "+err.Error()))
		return
	}

	observation, err := h.service.UpdateObservation(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to update observation")
		if err.Error() == "observation not found" {
			c.JSON(http.StatusNotFound, models.NewOperationOutcome("error", "not-found", "Observation not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewOperationOutcome("error", "exception", "Failed to update observation"))
		return
	}

	c.JSON(http.StatusOK, observation)
}

// DeleteObservation handles DELETE /api/v1/observations/:id
func (h *ObservationHandler) DeleteObservation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).WithField("id", idStr).Error("Invalid observation ID")
		c.JSON(http.StatusBadRequest, models.NewOperationOutcome("error", "invalid", "Invalid observation ID format"))
		return
	}

	err = h.service.DeleteObservation(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to delete observation")
		if err.Error() == "observation not found" {
			c.JSON(http.StatusNotFound, models.NewOperationOutcome("error", "not-found", "Observation not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewOperationOutcome("error", "exception", "Failed to delete observation"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListObservations handles GET /api/v1/observations
func (h *ObservationHandler) ListObservations(c *gin.Context) {
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

	response, err := h.service.ListObservations(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list observations")
		c.JSON(http.StatusInternalServerError, models.NewOperationOutcome("error", "exception", "Failed to list observations"))
		return
	}

	c.JSON(http.StatusOK, response)
}
