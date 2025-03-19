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

	router.POST("validate-google-id-token", authHandler.ValidateGoogleIdToken)
	router.POST("refresh-token", authHandler.RefreshToken)
	router.POST("auth/send-mail-test", authHandler.SendEmailVerificationExample)

	// Registration Route

	auth := router.Group("")
	auth.Use(authMiddleware.JwtAuthMiddleware)

	auth.GET("/auth/ping", authHandler.Ping)

	auth.GET("my-profile", userHandler.GetMyProfile)
	auth.POST("my-profile", userHandler.UpdateMyProfile)

	mustUpdatedProfile := auth.Group("")
	mustUpdatedProfile.Use(authMiddleware.MustUpdatedUserProfile)

	mustUpdatedProfile.GET("/authupdate/ping", authHandler.Ping)

	mustUpdatedProfile.POST("registration_team_hackathon", registrationHandler.Create)
	mustUpdatedProfile.POST("registration_user_hackathon", registrationHandler.UserJoinTeam)

	// Profile Route

	// router.POST("profile/post", profileHandler.Create)
	// router.GET("profile/get", profileHandler.GetProfile)

	// admin_api_base_url := os.Getenv("ADMIN_API_BASE_URL")
	// admin_router := r.Group(admin_api_base_url)
}
