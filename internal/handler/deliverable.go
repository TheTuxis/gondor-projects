package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheTuxis/gondor-projects/internal/model"
	"github.com/TheTuxis/gondor-projects/internal/service"
)

type DeliverableHandler struct {
	deliverableService *service.DeliverableService
}

func NewDeliverableHandler(deliverableService *service.DeliverableService) *DeliverableHandler {
	return &DeliverableHandler{deliverableService: deliverableService}
}

func (h *DeliverableHandler) List(c *gin.Context) {
	projectID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	deliverables, err := h.deliverableService.List(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to list deliverables",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deliverables})
}

func (h *DeliverableHandler) Create(c *gin.Context) {
	projectID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var input model.DeliverableCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "invalid request body",
			"details": gin.H{"error": err.Error()},
		})
		return
	}

	deliverable, err := h.deliverableService.Create(projectID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to create deliverable",
		})
		return
	}

	c.JSON(http.StatusCreated, deliverable)
}

func (h *DeliverableHandler) Update(c *gin.Context) {
	_, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	deliverableID, err := parseUintParam(c, "deliverable_id")
	if err != nil {
		return
	}

	var input model.DeliverableUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "invalid request body",
			"details": gin.H{"error": err.Error()},
		})
		return
	}

	deliverable, err := h.deliverableService.Update(deliverableID, input)
	if err != nil {
		if err == service.ErrDeliverableNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "deliverable not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to update deliverable",
		})
		return
	}

	c.JSON(http.StatusOK, deliverable)
}

func (h *DeliverableHandler) Delete(c *gin.Context) {
	_, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	deliverableID, err := parseUintParam(c, "deliverable_id")
	if err != nil {
		return
	}

	if err := h.deliverableService.Delete(deliverableID); err != nil {
		if err == service.ErrDeliverableNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "deliverable not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to delete deliverable",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deliverable deleted successfully"})
}
