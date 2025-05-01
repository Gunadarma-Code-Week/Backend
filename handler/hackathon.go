package handler

import (
	"gcw/dto"
	"gcw/helper"
	"gcw/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type hackathonHandler struct {
	service service.SubmissionService
}

type HackathonHandler interface {
	SubmissionHackaton(c *gin.Context)
	HackathonStageStatus(c *gin.Context)
}

func GateHackathonHandler(s service.SubmissionService) HackathonHandler {
	return &hackathonHandler{
		service: s,
	}
}

func (h *hackathonHandler) SubmissionHackaton(c *gin.Context) {
	stage := c.Param("stage")
	join_code := c.Param("join_code")

	var request dto.RequestHackathon
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", err.Error()))
		return
	}

	result, err := h.service.Create(join_code, stage, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("CREATE_FAILED", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateSuccessResponse("CREATED", result))
}

func (h *hackathonHandler) HackathonStageStatus(c *gin.Context) {
	join_code := c.Param("join_code")

	stageStatus, err := h.service.Get(join_code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("INTERNAL SERVER ERROR", "cant get stage status"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("SUCCESS", stageStatus))
}
