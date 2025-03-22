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
	response.User = *userResponse

	helper.SetTokenRefreshCookie(c, token, refreshToken)

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", response))
}

// @Summary Invalidate Cookie
// @Description Invalidate Cookie
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response
// @Router /auth/invalidate-cookie [post]
func (h *authHandler) InvalidateCookie(c *gin.Context) {
	helper.RemoveTokenRefreshCookie(c)
	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", nil))
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
