package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheTuxis/gondor-projects/internal/model"
	"github.com/TheTuxis/gondor-projects/internal/service"
)

type MemberHandler struct {
	memberService *service.MemberService
}

func NewMemberHandler(memberService *service.MemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

func (h *MemberHandler) List(c *gin.Context) {
	projectID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	members, err := h.memberService.List(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to list members",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": members})
}

func (h *MemberHandler) Create(c *gin.Context) {
	projectID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var input model.MemberCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation_error",
			"message": "invalid request body",
			"details": gin.H{"error": err.Error()},
		})
		return
	}

	member, err := h.memberService.Create(projectID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to add member",
		})
		return
	}

	c.JSON(http.StatusCreated, member)
}

func (h *MemberHandler) Delete(c *gin.Context) {
	_, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	memberID, err := parseUintParam(c, "member_id")
	if err != nil {
		return
	}

	if err := h.memberService.Delete(memberID); err != nil {
		if err == service.ErrMemberNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "member not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "failed to remove member",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed successfully"})
}
