package handler

import (
	"gcw/dto"
	"gcw/helper"
	"gcw/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TimelineHandler interface {
	CreateTimeline(c *gin.Context)
	UpdateTimeline(c *gin.Context)
	DeleteTimeline(c *gin.Context)
	GetTimelinesByCategory(c *gin.Context)
}

type timelineHandler struct {
	timelineService service.TimelineService
}

func NewTimelineHandler(ts service.TimelineService) TimelineHandler {
	return &timelineHandler{
		timelineService: ts,
	}
}

func formatError(err error, field string) interface{} {
	formatted := helper.FormatValidationError(err)
	if str, ok := formatted.(string); ok {
		return map[string][]string{field: {str}}
	}
	return formatted
}

func (h *timelineHandler) CreateTimeline(c *gin.Context) {
	var dtoCreate dto.TimelineCreateDTO
	err := c.ShouldBindJSON(&dtoCreate)
	if err != nil {
		res := helper.CreateErrorResponse("Failed to process request", formatError(err, "request"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := h.timelineService.CreateTimeline(dtoCreate)
	if err != nil {
		res := helper.CreateErrorResponse("Failed to create timeline", formatError(err, "timeline"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.CreateMutationResponse("Timeline created successfully", result)
	c.JSON(http.StatusCreated, res)
}

func (h *timelineHandler) UpdateTimeline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		res := helper.CreateErrorResponse("Invalid ID", map[string][]string{"id": {"IS_INVALID"}})
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var dtoUpdate dto.TimelineUpdateDTO
	err = c.ShouldBindJSON(&dtoUpdate)
	if err != nil {
		res := helper.CreateErrorResponse("Failed to process request", formatError(err, "request"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := h.timelineService.UpdateTimeline(uint(id), dtoUpdate)
	if err != nil {
		if err.Error() == "timeline not found" {
			c.JSON(http.StatusNotFound, helper.CreateNotFoundResponse("Timeline tidak ditemukan"))
			return
		}
		res := helper.CreateErrorResponse("Failed to update timeline", formatError(err, "timeline"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.CreateMutationResponse("Timeline updated successfully", result)
	c.JSON(http.StatusCreated, res)
}

func (h *timelineHandler) DeleteTimeline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		res := helper.CreateErrorResponse("Invalid ID", map[string][]string{"id": {"IS_INVALID"}})
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err = h.timelineService.DeleteTimeline(uint(id))
	if err != nil {
		if err.Error() == "timeline not found" {
			c.JSON(http.StatusNotFound, helper.CreateNotFoundResponse("Timeline tidak ditemukan"))
			return
		}
		res := helper.CreateErrorResponse("Failed to delete timeline", formatError(err, "timeline"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.CreateMutationResponse("Timeline deleted successfully", nil)
	c.JSON(http.StatusCreated, res)
}

func (h *timelineHandler) GetTimelinesByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		res := helper.CreateErrorResponse("Category is required", map[string][]string{"category": {"IS_REQUIRED"}})
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := h.timelineService.GetTimelinesByCategory(category)
	if err != nil {
		res := helper.CreateErrorResponse("Failed to fetch timelines", formatError(err, "timeline"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.CreateSuccessResponse("Timelines fetched successfully", result)
	c.JSON(http.StatusOK, res)
}
