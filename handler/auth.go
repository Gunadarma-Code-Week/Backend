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
	// response.AccessToken = token
	// response.RefreshToken = refreshToken

	// set cookie to client
	helper.SetTokenRefreshCookie(c, token, refreshToken)

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", response))
}

// regenerate the token
// func (h *authHandler) RefreshToken(c *gin.Context) {
// 	// refreshTokenDTO := &dto.RefreshTokenDTO{}
// 	// if err := c.Bind(refreshTokenDTO); err != nil {
// 	// 	logging.Low("AuthHandler.Login", "BAD_REQUEST", err.Error())
// 	// 	c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
// 	// 	return
// 	// }

// 	// get token from cookies
// 	oldRefreshToken, err := c.Cookie("gcw_api_refresh_token")
// 	if err != nil || oldRefreshToken == "" {
// 		c.JSON(401, helper.CreateErrorResponse("error", "refresh token is required"))
// 		return
// 	}

// 	claims, err := h.jwtService.GetClaimsByRefreshToken(oldRefreshToken)
// 	if err != nil {
// 		logging.High("AuthHandler.Login", "INVALID_TOKEN", err.Error())
// 		c.JSON(401, helper.CreateErrorResponse("error", "invalid token"))
// 		return
// 	}

// 	email, ok := claims["email"].(string)
// 	if !ok {
// 		c.JSON(401, helper.CreateErrorResponse("error", "email not found in token"))
// 		return
// 	}

// 	user, err := h.authService.FindByEmail(email)
// 	if err != nil {
// 		logging.High("AuthHandler.Login", "USER_NOT_FOUND", err.Error())
// 		c.JSON(401, helper.CreateErrorResponse("error", "user not found"))
// 		return
// 	}

// 	token := h.jwtService.GenerateToken(user)
// 	refreshToken := h.jwtService.GenerateRefreshToken(user)

// 	userResponse := &dto.UserResponseDTO{}
// 	smapping.FillStruct(userResponse, smapping.MapFields(user))

// 	response := &dto.AuthResponseDTO{}
// 	response.User = *userResponse
// 	response.AccessToken = token
// 	response.RefreshToken = refreshToken

// 	c.SetCookie("gcw_api_token", token, 60*60*12, "/", "localhost", false, true)
// 	c.SetCookie("gcw_api_refresh_token", refreshToken, 60*60*24*5, "/", "localhost", false, true)

// 	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", response))
// }

// THIS JUST EXAMPLE, CAN USE THIS ON ANYWHERE
func (h *authHandler) SendEmailVerificationExample(c *gin.Context) {
	// use gorooutine to send email, so it will not blocking the main process
	// u can use goroutine on any process that not need to wait the process
	go h.emailService.SendEmailHTML("Email Verification", []string{"tes@mail.com"}, "template/email/verification.html", map[string]string{
		"Code": "123456",
	})

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", "Email verification has been sent, wait or try again"))
}
