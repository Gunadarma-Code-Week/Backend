package handler

import (
	"gcw/dto"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"log"
	"net/http"
	"strings"

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
	c.JSON(http.StatusOK, helper.CreateSuccessResponse("ping", "ping"))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	user, err := h.authService.GetUserByGoogleIdToken(login.GoogleIdToken)
	if err != nil {
		logging.High("AuthHandler.Login", "AUTH_ERROR", err.Error())
		
		statuscode := http.StatusBadRequest
		var errorPayload interface{}
		if strings.Contains(err.Error(), "dinonaktifkan") {
			statuscode = http.StatusBadRequest
			errorPayload = map[string][]string{"account": {"IS_INVALID"}}
		} else {
			errorPayload = map[string][]string{"credential": {"IS_INVALID"}}
		}
		c.JSON(statuscode, helper.CreateErrorResponse(err.Error(), errorPayload))
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

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Data berhasil diperbarui", response))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	payload, err := h.jwtService.GetClaimsByRefreshToken(refreshToken.RefreshToken)
	if err != nil {
		logging.High("AuthHandler.RefreshToken", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	userId := uint64(payload["id"].(float64))
	if userId == 0 {
		logging.High("AuthHandler.RefreshToken", "INTERNAL_SERVER_ERROR", "user_id not found in payload")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"user_id": {"IS_INVALID"}}))
		return
	}

	user, err := h.authService.GetUserById(userId)
	if err != nil {
		logging.High("AuthHandler.RefreshToken", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
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

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Data berhasil diperbarui", response))
}

// @Summary Register new account
// @Description Register a new user account with email, password, and name
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body dto.RegisterDTO true "Register Payload"
// @Success 200 {object} helper.Response{data=dto.AuthResponseDTO}
// @Failure 400 {object} helper.Response
// @Router /auth/registration [post]
func (h *authHandler) Registration(c *gin.Context) {
	// Get the validated DTO from the context
	auth, _ := c.Get("dto")
	registerDTO := auth.(*dto.RegisterDTO)

	user, err := h.authService.Registration(registerDTO)
	if err != nil {
		if err.Error() == "user already registered" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("Email sudah terdaftar"))
			return
		}
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
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

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Data berhasil diperbarui", response))
}

// @Summary Login
// @Description Login with email and password to get access token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body dto.LoginDTO true "Login Payload"
// @Success 200 {object} helper.Response{data=dto.AuthResponseDTO}
// @Failure 400 {object} helper.Response
// @Router /auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	// Get the validated DTO from the context
	auth, _ := c.Get("dto")
	loginDTO := auth.(*dto.LoginDTO)

	user, err := h.authService.LoginService(loginDTO)
	if err != nil {
		var errorPayload interface{}
		if strings.Contains(err.Error(), "dinonaktifkan") {
			errorPayload = map[string][]string{"account": {"IS_INVALID"}}
		} else {
			errorPayload = map[string][]string{"credential": {"IS_INVALID"}}
		}
		msg := err.Error()
		if msg == "email or password is wrong" {
			msg = "Email atau password salah"
		}
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse(msg, errorPayload))
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

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Data berhasil diperbarui", response))
}

// THIS JUST EXAMPLE, CAN USE THIS ON ANYWHERE
func (h *authHandler) SendEmailVerificationExample(c *gin.Context) {
	// use gorooutine to send email, so it will not blocking the main process
	// u can use goroutine on any process that not need to wait the process
	go h.emailService.SendEmailHTML("Email Verification", []string{"tes@mail.com"}, "template/email/verification.html", map[string]string{
		"Code": "123456",
	})

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Data berhasil diperbarui", "Email verification has been sent, wait or try again"))
}
