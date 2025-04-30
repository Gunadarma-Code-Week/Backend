package handler

import (
	"gcw/dto"
	"gcw/entity"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

type registrationHandler struct {
	registrationService *service.RegistrationService
	userService         *service.UserService
}

func GateRegistrationHandler(service *service.RegistrationService, userService *service.UserService) *registrationHandler {
	return &registrationHandler{
		registrationService: service,
		userService:         userService,
	}
}

// @Summary Register CP Team
// @Tags Team Registration
// @Accept json
// @Produce  json
// @Param request body dto.RegistrationCPTeamRequest true "Register CP Team"
// @Success 200 {object} helper.Response{data=dto.RegistrationCPTeamResponse}
// @Router /team/registration/cp [post]
func (h *registrationHandler) RegistrationCPTeam(c *gin.Context) {
	registrationDto := &dto.RegistrationCPTeamRequest{}

	if err := c.ShouldBind(registrationDto); err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	userAuth := c.MustGet("user").(*entity.User)

	registrationCPTeamResponse, err := h.registrationService.CPTeamRegistration(registrationDto, userAuth)
	if err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateSuccessResponse("Success register cp team", registrationCPTeamResponse))
}

// @Summary Register Hackathon Team
// @Tags Team Registration
// @Accept json
// @Produce  json
// @Param request body dto.RegistrationHackathonTeamRequest true "Register Hackathon Team"
// @Success 200 {object} helper.Response{data=dto.RegistrationHackathonTeamResponse}
// @Router /team/registration/hackathon [post]
func (h *registrationHandler) RegistrationHackathonTeam(c *gin.Context) {
	registrationDto := &dto.RegistrationHackathonTeamRequest{}

	if err := c.ShouldBind(registrationDto); err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	userAuth := c.MustGet("user").(*entity.User)

	registrationHackathonTeamResponse, err := h.registrationService.HackathonTeamRegistration(registrationDto, userAuth)
	if err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateSuccessResponse("Success register hackathon team", registrationHackathonTeamResponse))
}

// @Summary Find Team
// @Tags Team Registration
// @Produce  json
// @Param join_code query string true "Join Code"
// @Success 200 {object} helper.Response{data=dto.RegistraionTeamResponse}
// @Router /team/registration/find/{join_code} [get]
func (h *registrationHandler) FindTeam(c *gin.Context) {
	joinCode := c.Param("join_code")

	team, err := h.registrationService.FindTeamByJoinCode(joinCode)
	if err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	registraionTeamResponse := &dto.RegistraionTeamResponse{}
	err = smapping.FillStruct(registraionTeamResponse, smapping.MapFields(team))
	if err != nil {
		logging.Low("RegistrationService.CPTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	members, err := h.userService.FindByIdTeam(team.ID_Team)
	if err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	registraionTeamResponse.Members = members

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Success Find Team", registraionTeamResponse))
}

// @Summary Join Team
// @Tags Team Registration
// @Produce  json
// @Param join_code query string true "Join Code"
// @Success 200 {object} helper.Response{data=dto.RegistraionTeamResponse}
// @Router /team/registration/join/{join_code} [post]
func (h *registrationHandler) UserJoinTeam(c *gin.Context) {
	joinCode := c.Query("join_code")

	userAuth := c.MustGet("user").(*entity.User)

	team, err := h.registrationService.JoinTeam(joinCode, userAuth)
	if err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	registraionTeamResponse := &dto.RegistraionTeamResponse{}
	err = smapping.FillStruct(registraionTeamResponse, smapping.MapFields(team))
	if err != nil {
		logging.Low("RegistrationService.CPTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Success Join Team", registraionTeamResponse))
}
