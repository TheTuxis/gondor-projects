package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheTuxis/gondor-projects/internal/model"
	"github.com/TheTuxis/gondor-projects/internal/service"
)

type PhaseHandler struct {
	phaseService *service.PhaseService
}

func NewPhaseHandler(phaseService *service.PhaseService) *PhaseHandler {
	return &PhaseHandler{phaseService: phaseService}
}

func (h *PhaseHandler) List(c *gin.Context) {
	projectID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	phases, err := h.phaseService.List(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to list phases",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": phases})
}

func (h *PhaseHandler) Create(c *gin.Context) {
	projectID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var input model.PhaseCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "invalid request body",
			"details": gin.H{"error": err.Error()},
		})
		return
	}

	phase, err := h.phaseService.Create(projectID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to create phase",
		})
		return
	}

	c.JSON(http.StatusCreated, phase)
}

func (h *PhaseHandler) Update(c *gin.Context) {
	_, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	phaseID, err := parseUintParam(c, "phase_id")
	if err != nil {
		return
	}

	var input model.PhaseUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "invalid request body",
			"details": gin.H{"error": err.Error()},
		})
		return
	}

	phase, err := h.phaseService.Update(phaseID, input)
	if err != nil {
		if err == service.ErrPhaseNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "phase not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to update phase",
		})
		return
	}

	c.JSON(http.StatusOK, phase)
}

func (h *PhaseHandler) Delete(c *gin.Context) {
	_, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	phaseID, err := parseUintParam(c, "phase_id")
	if err != nil {
		return
	}

	if err := h.phaseService.Delete(phaseID); err != nil {
		if err == service.ErrPhaseNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "phase not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to delete phase",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "phase deleted successfully"})
}
