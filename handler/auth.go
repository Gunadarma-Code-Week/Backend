package handler

import (
	"gcw/dto"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

type authHandler struct {
	authService  *service.AuthService
	jwtService   *service.JwtService
	emailService *service.EmailService
}

func NewAuthHandler(as *service.AuthService, js *service.JwtService, es *service.EmailService) *authHandler {
	return &authHandler{
		authService:  as,
		jwtService:   js,
		emailService: es,
	}
}

func (h *authHandler) Ping(c *gin.Context) {
	log.Printf("Berhasil Ping")
	c.JSON(http.StatusOK, gin.H{"success": "ping"})
}

// @Summary Validate Google ID Token
// @Description Validate Google ID Token (Login)
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body dto.ValidateGoogleIdTokenDTO true "Google ID Token"
// @Success 200 {object} helper.Response{data=dto.AuthResponseDTO}
// @Router /auth/validate-google-id-token [post]
func (h *authHandler) ValidateGoogleIdToken(c *gin.Context) {
	login := &dto.ValidateGoogleIdTokenDTO{}
	if err := c.Bind(login); err != nil {
		logging.Low("AuthHandler.Login", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	user, err := h.authService.GetUserByGoogleIdToken(login.GoogleIdToken)
	if err != nil {
		logging.High("AuthHandler.Login", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	token := h.jwtService.GenerateToken(user)
	refreshToken := h.jwtService.GenerateRefreshToken(user)

	userResponse := &dto.UserResponseDTO{}
	smapping.FillStruct(userResponse, smapping.MapFields(user))

	response := &dto.AuthResponseDTO{}
	response.AccessToken = token
	response.RefreshToken = refreshToken
	response.User = *userResponse

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", response))
}

// @Summary Refresh Token
// @Description Refresh Token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body dto.RefreshTokenDTO true "Refresh Token"
// @Success 200 {object} helper.Response{data=dto.AuthResponseDTO}
// @Router /auth/refresh-token [post]
func (h *authHandler) RefreshToken(c *gin.Context) {
	refreshToken := &dto.RefreshTokenDTO{}
	if err := c.Bind(refreshToken); err != nil {
		logging.Low("AuthHandler.RefreshToken", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	payload, err := h.jwtService.GetClaimsByRefreshToken(refreshToken.RefreshToken)
	if err != nil {
		logging.High("AuthHandler.RefreshToken", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	userId := uint64(payload["id"].(float64))
	if userId == 0 {
		logging.High("AuthHandler.RefreshToken", "INTERNAL_SERVER_ERROR", "user_id not found in payload")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "user_id not found in payload"))
		return
	}

	user, err := h.authService.GetUserById(userId)
	if err != nil {
		logging.High("AuthHandler.RefreshToken", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}
	newToken := h.jwtService.GenerateToken(user)
	newRefreshToken := h.jwtService.GenerateRefreshToken(user)

	userResponse := &dto.UserResponseDTO{}
	smapping.FillStruct(userResponse, smapping.MapFields(user))

	response := &dto.AuthResponseDTO{}
	response.AccessToken = newToken
	response.RefreshToken = newRefreshToken
	response.User = *userResponse

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", response))
}

func (h *authHandler) Registration(c *gin.Context) {
	auth := &dto.AuthenticationDTO{}
	user, err := h.authService.Registration(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", err.Error()))
		return
	}

	token := h.jwtService.GenerateToken(user)
	refreshToken := h.jwtService.GenerateRefreshToken(user)

	userResponse := &dto.UserResponseDTO{}
	smapping.FillStruct(userResponse, smapping.MapFields(user))

	response := &dto.AuthResponseDTO{}
	response.AccessToken = token
	response.RefreshToken = refreshToken
	response.User = *userResponse

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", response))
}

func (h *authHandler) Login(c *gin.Context) {
	auth := &dto.AuthenticationDTO{}
	user, err := h.authService.LoginService(auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", err.Error()))
		return
	}

	token := h.jwtService.GenerateToken(user)
	refreshToken := h.jwtService.GenerateRefreshToken(user)

	userResponse := &dto.UserResponseDTO{}
	smapping.FillStruct(userResponse, smapping.MapFields(user))

	response := &dto.AuthResponseDTO{}
	response.AccessToken = token
	response.RefreshToken = refreshToken
	response.User = *userResponse

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", response))
}

// THIS JUST EXAMPLE, CAN USE THIS ON ANYWHERE
func (h *authHandler) SendEmailVerificationExample(c *gin.Context) {
	// use gorooutine to send email, so it will not blocking the main process
	// u can use goroutine on any process that not need to wait the process
	go h.emailService.SendEmailHTML("Email Verification", []string{"tes@mail.com"}, "template/email/verification.html", map[string]string{
		"Code": "123456",
	})

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", "Email verification has been sent, wait or try again"))
}
