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

func (h *timelineHandler) CreateTimeline(c *gin.Context) {
	var dtoCreate dto.TimelineCreateDTO
	err := c.ShouldBindJSON(&dtoCreate)
	if err != nil {
		res := helper.CreateErrorResponse("Failed to process request", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := h.timelineService.CreateTimeline(dtoCreate)
	if err != nil {
		res := helper.CreateErrorResponse("Failed to create timeline", err.Error())
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
		res := helper.CreateErrorResponse("Invalid ID", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var dtoUpdate dto.TimelineUpdateDTO
	err = c.ShouldBindJSON(&dtoUpdate)
	if err != nil {
		res := helper.CreateErrorResponse("Failed to process request", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := h.timelineService.UpdateTimeline(uint(id), dtoUpdate)
	if err != nil {
		res := helper.CreateErrorResponse("Failed to update timeline", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.CreateMutationResponse("Timeline updated successfully", result)
	c.JSON(http.StatusOK, res)
}

func (h *timelineHandler) DeleteTimeline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		res := helper.CreateErrorResponse("Invalid ID", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err = h.timelineService.DeleteTimeline(uint(id))
	if err != nil {
		res := helper.CreateErrorResponse("Failed to delete timeline", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.CreateMutationResponse("Timeline deleted successfully", nil)
	c.JSON(http.StatusOK, res)
}

func (h *timelineHandler) GetTimelinesByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		res := helper.CreateErrorResponse("Category is required", "Empty category")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := h.timelineService.GetTimelinesByCategory(category)
	if err != nil {
		res := helper.CreateErrorResponse("Failed to fetch timelines", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.CreateSuccessResponse("Timelines fetched successfully", result)
	c.JSON(http.StatusOK, res)
}
