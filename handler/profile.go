package handler

import (
	"gcw/dto"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

type profileHandler struct {
	profileService service.ProfileService
}

type ProfileHandler interface {
	getProfile(*gin.Context)
}

func GateProfileHandler(service service.ProfileService) ProfileHandler {
	return &profileHandler{
		profileService: service,
	}
}

func (h *profileHandler) getProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id_user"), 0, 0)

	if err != nil {
		logging.Low("ProfileHandler.Register", "NOT_FOUND", err.Error())
		c.JSON(http.StatusNotFound, helper.CreateErrorResponse("Data tidak ditemukan", "Not Found"))
		return
	}

	data, err := h.profileService.Get(id)

	response := &dto.ProfileResponseDTO{}
	smapping.FillStruct(response, smapping.MapFields(data))

	if err != nil {
		logging.Low("ProfileHandler.getProfile", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusNoContent, helper.CreateSuccessResponse("Data successfully obtained", response))
}

func (h *profileHandler) Create(c *gin.Context) {
	profileDto := &dto.ProfileResponseDTO{}

	if err := c.ShouldBind(profileDto); err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	profile, err := h.profileService.Create(profileDto)

	response := &dto.ProfileResponseDTO{}
	smapping.FillStruct(response, smapping.MapFields(profile))

	if err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateSuccessResponse("Success create data", response))
}
