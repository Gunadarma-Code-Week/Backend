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

// @Summary Submit Hackathon
// @Description Submit a hackathon entry
// @Tags Hackathon
// @Accept  json
// @Produce  json
// @Param stage path string true "Hackathon Stage"
// @Param join_code path string true "Hackathon Join Code"
// @Param request body dto.RequestHackathon true "Hackathon Submission Request"
// @Success 201 {object} helper.Response{data=entity.HackathonTeam}
// @Failure 400 {object} helper.Response{data=string} "Invalid input"
// @Failure 500 {object} helper.Response{data=string} "Submission failed"
// @Router /hackathon/{join_code}/{stage} [post]
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

// @Summary Get Hackathon Stage Status
// @Description Get the status of a hackathon stage
// @Tags Hackathon
// @Accept  json
// @Produce  json
// @Param join_code path string true "Hackathon Join Code"
// @Success 200 {object} helper.Response{data=dto.HackatonStageStatus}
// @Failure 400 {object} helper.Response{data=string} "Invalid join code"
// @Router /hackathon/{join_code}/status [get]
func (h *hackathonHandler) HackathonStageStatus(c *gin.Context) {
	join_code := c.Param("join_code")

	stageStatus, err := h.service.Get(join_code)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("SUCCESS", stageStatus))
}
