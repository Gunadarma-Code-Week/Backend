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

	userRepository = repository.NewUserRepository(database)
	// profileRepository      = repository.GateProfileRepository(database)
	registrationRepository = repository.GateRegistrationRepository(database)

	authService         = service.NewAuthService(userRepository)
	userService         = service.NewUserService(userRepository)
	registrationService = service.GatRegistrationService(registrationRepository)
	jwtService          = service.NewJwtService()
	emailService        = service.NewEmailService()

	authHandler         = handler.NewAuthHandler(authService, jwtService, emailService)
	userHandler         = handler.NewUserHandler(userService)
	registrationHandler = handler.GateRegistrationHandler(registrationService)

	authMiddleware = middleware.NewAuthMiddleware(authService, jwtService)
)

func SetupRouter(r *gin.Engine) {
	api_base_url := os.Getenv("API_BASE_URL")
	router := r.Group(api_base_url)

	router.GET("/ping", authHandler.Ping)

	// Authentication Route
	auth := router.Group("auth")
	auth.POST("validate-google-id-token", authHandler.ValidateGoogleIdToken)
	auth.POST("send-mail-test", authHandler.SendEmailVerificationExample)

	mustAuth := router.Group("")
	mustAuth.Use(authMiddleware.JwtAuthMiddleware)
	mustAuth.GET("mustauth/ping", authHandler.Ping)

	// Profile Route
	profile := mustAuth.Group("profile")
	profile.GET("my", userHandler.GetMyProfile)
	profile.POST("my", userHandler.UpdateMyProfile)

	mustUpdatedProfile := mustAuth.Group("")
	mustUpdatedProfile.Use(authMiddleware.MustUpdatedUserProfile)
	mustUpdatedProfile.GET("mustauth/authupdate/ping", authHandler.Ping)

	// Team Registrasion
	teamRegistration := mustUpdatedProfile.Group("team/registration")
	teamRegistration.POST("hackathon", registrationHandler.Create)
	teamRegistration.POST("hackathon/join", registrationHandler.UserJoinTeam)

	// Profile Route

	// router.POST("profile/post", profileHandler.Create)
	// router.GET("profile/get", profileHandler.GetProfile)

	// admin_api_base_url := os.Getenv("ADMIN_API_BASE_URL")
	// admin_router := r.Group(admin_api_base_url)
}
