package handler

import (
	"fmt"
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	userAuth := c.MustGet("user").(*entity.User)

	registrationCPTeamResponse, err := h.registrationService.CPTeamRegistration(registrationDto, userAuth)
	if err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		if err.Error() == "Nama Tim Sudah Digunakan" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("Nama tim sudah digunakan"))
			return
		} else if err.Error() == "USER ALREADY HAVE TEAM" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("User sudah memiliki tim"))
			return
		} else if err.Error() == "Pendaftaran Competitive Programming telah ditutup" {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse(err.Error(), map[string][]string{"registration": {"IS_INVALID"}}))
			return
		}
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	c.Set("target_name", registrationCPTeamResponse.Team.TeamName)
	c.Set("target_id", uint64(registrationCPTeamResponse.Team.ID_Team))
	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Success register cp team", registrationCPTeamResponse))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	userAuth := c.MustGet("user").(*entity.User)

	registrationHackathonTeamResponse, err := h.registrationService.HackathonTeamRegistration(registrationDto, userAuth)
	if err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		if err.Error() == "Nama Tim Sudah Digunakan" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("Nama tim sudah digunakan"))
			return
		} else if err.Error() == "USER ALREADY HAVE TEAM" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("User sudah memiliki tim"))
			return
		} else if err.Error() == "Pendaftaran Hackathon telah ditutup" {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse(err.Error(), map[string][]string{"registration": {"IS_INVALID"}}))
			return
		}
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	c.Set("target_name", registrationHackathonTeamResponse.Team.TeamName)
	c.Set("target_id", uint64(registrationHackathonTeamResponse.Team.ID_Team))
	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Success register hackathon team", registrationHackathonTeamResponse))
}

// @Summary Register CTF Team
// @Tags Team Registration
// @Accept json
// @Produce  json
// @Param request body dto.RegistrationCTFTeamRequest true "Register CTF Team"
// @Success 200 {object} helper.Response{data=dto.RegistrationCTFTeamResponse}
// @Router /team/registration/ctf [post]
func (h *registrationHandler) RegistrationCTFTeam(c *gin.Context) {
	registrationDto := &dto.RegistrationCTFTeamRequest{}

	if err := c.ShouldBind(registrationDto); err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	userAuth := c.MustGet("user").(*entity.User)

	registrationCTFTeamResponse, err := h.registrationService.CTFTeamRegistration(registrationDto, userAuth)
	if err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		if err.Error() == "Nama Tim Sudah Digunakan" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("Nama tim sudah digunakan"))
			return
		} else if err.Error() == "USER ALREADY HAVE TEAM" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("User sudah memiliki tim"))
			return
		} else if err.Error() == "Pendaftaran Capture The Flag telah ditutup" {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse(err.Error(), map[string][]string{"registration": {"IS_INVALID"}}))
			return
		}
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	c.Set("target_name", registrationCTFTeamResponse.Team.TeamName)
	c.Set("target_id", uint64(registrationCTFTeamResponse.Team.ID_Team))
	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Success register ctf team", registrationCTFTeamResponse))
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
		logging.Low("RegistrationHandler.FindTeam", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	registrationTeamResponse := &dto.RegistraionTeamResponse{}
	err = smapping.FillStruct(registrationTeamResponse, smapping.MapFields(team))
	if err != nil {
		logging.Low("RegistrationHandler.FindTeam", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	leader, err := h.userService.FindById(team.ID_LeadTeam)
	if err != nil {
		logging.Low("RegistrationHandler.FindTeam", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	members, err := h.userService.FindByIdTeam(team.ID_Team, leader.ID)
	if err != nil {
		logging.Low("RegistrationHandler.FindTeam", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	leaderData := dto.Member{
		Name:       leader.Name,
		Email:      leader.Email,
		University: leader.Institusi,
		Role:       "Leader",
	}

	fmt.Println(leaderData)

	registrationTeamResponse.Members = members
	registrationTeamResponse.Leader = leaderData

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Success Find Team", registrationTeamResponse))
}

// @Summary Join Team
// @Tags Team Registration
// @Produce  json
// @Param join_code query string true "Join Code"
// @Success 200 {object} helper.Response{data=dto.RegistraionTeamResponse}
// @Router /team/registration/join/{join_code} [post]
func (h *registrationHandler) UserJoinTeam(c *gin.Context) {
	joinCode := c.Param("join_code")

	userAuth := c.MustGet("user").(*entity.User)

	team, err := h.registrationService.JoinTeam(joinCode, userAuth)
	if err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		if err.Error() == "USER ALREADY HAVE TEAM" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("User sudah memiliki tim"))
			return
		} else if err.Error() == "TEAM IS FULL" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("Tim sudah penuh"))
			return
		} else if err.Error() == "TEAM NOT FOUND" {
			c.JSON(http.StatusNotFound, helper.CreateNotFoundResponse("Tim tidak ditemukan"))
			return
		} else if err.Error() == "Pendaftaran Hackathon telah ditutup" || err.Error() == "Pendaftaran Competitive Programming telah ditutup" || err.Error() == "Pendaftaran Capture The Flag telah ditutup" {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse(err.Error(), map[string][]string{"registration": {"IS_INVALID"}}))
			return
		}
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	registraionTeamResponse := &dto.RegistraionTeamResponse{}
	err = smapping.FillStruct(registraionTeamResponse, smapping.MapFields(team))
	if err != nil {
		logging.Low("RegistrationService.CPTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	c.Set("target_name", team.TeamName)
	c.Set("target_id", uint64(team.ID_Team))
	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Success Join Team", registraionTeamResponse))
}
