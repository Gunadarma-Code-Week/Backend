package router

import (
	"gcw/config"
	"gcw/handler"
	"gcw/middleware"
	"gcw/repository"
	"gcw/service"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	database = config.SetupDatabaseConnection()

	userRepository         = repository.NewUserRepository(database)
	profileRepository      = repository.GateProfileRepository(database)
	registrationRepository = repository.GateRegistrationRepository(database)

	authService         = service.NewAuthService(userRepository)
	profileService      = service.GateProfileService(profileRepository)
	registrationService = service.GatRegistrationService(registrationRepository)
	jwtService          = service.NewJwtService()
	emailService        = service.NewEmailService()

	authHandler         = handler.NewAuthHandler(authService, jwtService, emailService)
	profileHandler      = handler.GateProfileHandler(profileService)
	registrationHandler = handler.GateRegistrationHandler(registrationService)

	authMiddleware = middleware.NewAuthMiddleware(authService, jwtService)
)

func SetupRouter(r *gin.Engine) {
	api_base_url := os.Getenv("API_BASE_URL")
	router := r.Group(api_base_url)

	router.GET("/ping", authHandler.Ping)

	// Authentication Route

	router.POST("validate-google-id-token", authHandler.ValidateGoogleIdToken)
	router.POST("auth/send-mail-test", authHandler.SendEmailVerificationExample)

	// Registration Route

	auth := router.Group("auth")
	auth.Use(authMiddleware.JwtAuthMiddleware)

	auth.GET("/ping", authHandler.Ping)

	auth.POST("registration_team_hackathon", registrationHandler.Create)
	auth.POST("registration_user_hackathon", registrationHandler.UserJoinTeam)

	// Profile Route

	router.POST("profile/post", profileHandler.Create)
	router.GET("profile/get", profileHandler.GetProfile)

	// admin_api_base_url := os.Getenv("ADMIN_API_BASE_URL")
	// admin_router := r.Group(admin_api_base_url)
}
