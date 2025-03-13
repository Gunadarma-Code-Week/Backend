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
	authService  service.AuthService
	jwtService   service.JwtService
	emailService *service.EmailService
}

// UNECESSARY USING INTERFACE IF JUST ONE IMPLEMENTATION, JUST USE STRUCT
// u can remove entire interface, n reduce using pointer on unecessary usage just pass by value
// it can reduce some memory usage / overhead

// type AuthHandler interface {
// 	Ping(*gin.Context)
// 	Register(*gin.Context)
// 	Login(*gin.Context)
// }

func NewAuthHandler(as service.AuthService, js service.JwtService, es *service.EmailService) *authHandler {
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

func (h *authHandler) Login(c *gin.Context) {
	login := &dto.LoginDTO{}
	if err := c.Bind(login); err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}
	user, err := h.authService.FindByUsername(login.Username)

	if err != nil {
		logging.Low("AuthHandler.Login", "INTERNAL_SERVER_ERROR", "Username Not Found")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Username tidak ditemukan", err.Error()))
		return
	}

	if res := (h.authService.VerifyPassword(user.Password, login.Password)); !res {
		logging.Low("AuthHandler.Login", "INTERNAL_SERVER_ERROR", "Wrong Password")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Pssword salah", "wrong password"))
		return
	}

	token := h.jwtService.GenerateToken(user.Username)

	response := &dto.UserResponseDTO{}
	smapping.FillStruct(response, smapping.MapFields(user))
	response.AccessToken = token

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", response))
}

func (h *authHandler) Register(c *gin.Context) {
	register := &dto.UserRequestDTO{}

	if err := c.Bind(register); err != nil {
		logging.Low("AuthHandler.Register", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	print(register)

	user, err := h.authService.Register(register)

	if err != nil {
		logging.High("AuthHandler.Register", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	token := h.jwtService.GenerateToken(user.Username)

	response := &dto.UserResponseDTO{}
	smapping.FillStruct(response, smapping.MapFields(user))
	response.AccessToken = token

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
