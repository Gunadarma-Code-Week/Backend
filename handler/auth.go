package handler

import (
	"gcw/dto"
	"gcw/helper"
	"gcw/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

type authHandler struct {
	authService service.AuthService
	jwtService  service.JwtService
}

type AuthHandler interface {
	Ping(*gin.Context)
	Register(*gin.Context)
}

func NewAuthHandler(as service.AuthService, js service.JwtService) AuthHandler {
	return &authHandler{
		authService: as,
		jwtService:  js,
	}
}

func (h *authHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": "ping"})
}

func (h *authHandler) Register(c *gin.Context) {
	register := &dto.UserRequestDTO{}
	if err := c.Bind(register); err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	user, err := h.authService.Register(register)

	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	token := h.jwtService.GenerateToken(user.Username)

	response := &dto.UserResponseDTO{}
	smapping.FillStruct(response, smapping.MapFields(user))
	response.AccessToken = token

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", response))
}
