package router

import (
	"gcw/config"
	"gcw/handler"
	"gcw/repository"
	"gcw/service"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	database = config.SetupDatabaseConnection()

	userRepository    = repository.NewUserRepository(database)
	profileRepository = repository.GateProfileRepository(database)

	authService    = service.NewAuthService(userRepository)
	profileService = service.GateProfileService(profileRepository)
	jwtService     = service.NewJwtService()
	emailService   = service.NewEmailService()

	authHandler    = handler.NewAuthHandler(authService, jwtService, emailService)
	profileHandler = handler.GateProfileHandler(profileService)
)

func SetupRouter(r *gin.Engine) {
	api_base_url := os.Getenv("API_BASE_URL")
	router := r.Group(api_base_url)

	router.GET("/ping", authHandler.Ping)

	// Authentication Route

	router.POST("login", authHandler.Login)
	router.POST("register", authHandler.Register)
	router.POST("auth/send-mail-test", authHandler.SendEmailVerificationExample)

	// router.POST("profile/post", profileHandler.Create)
	// router.GET("profile/get", profileHandler.GetProfile)

	// admin_api_base_url := os.Getenv("ADMIN_API_BASE_URL")
	// admin_router := r.Group(admin_api_base_url)
}
