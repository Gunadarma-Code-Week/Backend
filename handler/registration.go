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

type RegistrationHandler interface {
	Create(*gin.Context)
}

func GateRegistrationHandler(service service.RegistrationService) RegistrationHandler {
	return &registrationHandler{
		registrationService: service,
	}
}

func (h *registrationHandler) Create(c *gin.Context) {
	registrationDto := &dto.RegistrationResponseDTO{}

	if err := c.ShouldBind(registrationDto); err != nil {
		logging.Low("ProfileHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	registration, err := h.registrationService.Create(registrationDto)

	response := &dto.RegistrationResponseDTO{}
	smapping.FillStruct(response, smapping.MapFields(registration))

	if err != nil {
		logging.Low("RegistrationHandler.Create", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateSuccessResponse("Success create data", response))
}
