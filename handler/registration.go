package handler

import (
	"gcw/dto"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

type registrationHandler struct {
	registrationService service.RegistrationService
}

// type RegistrationHandler interface {
// 	Create(*gin.Context)
// }

func GateRegistrationHandler(service service.RegistrationService) *registrationHandler {
	return &registrationHandler{
		registrationService: service,
	}
}

func (h *registrationHandler) Create(c *gin.Context) {
	registrationDto := &dto.RegistrationRequestHackathon{}

	if err := c.ShouldBind(registrationDto); err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	joinCode := helper.GenerateJoinCode()
	registrationDto.JoinCode = joinCode

	team := registrationDto.RegistrationResponseWithJoinCode
	hackathonTeam := registrationDto.RegistrationResponseHackathon

	registration, err := h.registrationService.Create(&team)

	if err != nil {
		logging.Low("RegistrationHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	id_team := registration.ID_Team

	hackathonTeam.IDTeam = id_team

	if _, err := h.registrationService.CreateTeam(&hackathonTeam); err != nil {
		logging.Low("RegistrationHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	if _, err := h.registrationService.UpdateUser(id_team, team.ID_LeadTeam); err != nil {
		logging.Low("RegistrationHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	combinedResponse := dto.RegistrationCombinedResponse{
		Registration:  team,
		HackathonTeam: hackathonTeam,
	}

	c.JSON(http.StatusCreated, helper.CreateSuccessResponse("Success create data", combinedResponse))
}

func (h *registrationHandler) UserJoinTeam(c *gin.Context) {
	registrationDto := &dto.JoinTeam{}

	if err := c.ShouldBind(registrationDto); err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	_, err := h.registrationService.UpdateUserJoinCode(registrationDto.TeamCode, registrationDto.IDUser)

	response := &dto.JoinTeam{}
	smapping.FillStruct(response, smapping.MapFields(registrationDto))

	if err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateSuccessResponse("Success create data", response))
}
