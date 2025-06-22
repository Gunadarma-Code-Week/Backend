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

	userRepository       = repository.NewUserRepository(database)
	newsletterRepository = repository.NewNewsletterRepository(database)
	// profileRepository      = repository.GateProfileRepository(database)
	registrationRepository = repository.GateRegistrationRepository(database)

	jwtService          = service.NewJwtService()
	emailService        = service.NewEmailService()
	domJudgeService     = service.NewDomJudgeService()
	authService         = service.NewAuthService(userRepository)
	userService         = service.NewUserService(userRepository)
	registrationService = service.NewRegistrationService(
		registrationRepository,
		domJudgeService,
	)
	newsletterService = service.NewNewsletterService(newsletterRepository)
	SubmissionService = service.NewSubmissionService(database)
	cpService         = service.NewCpService(database)

	authHandler         = handler.NewAuthHandler(authService, jwtService, emailService)
	userHandler         = handler.NewUserHandler(userService)
	registrationHandler = handler.GateRegistrationHandler(registrationService, userService)
	newsletterHandler   = handler.NewNewsletterHandler(newsletterService)
	submissionHandler   = handler.GateHackathonHandler(SubmissionService)
	cpHandler           = handler.GateCompetitiveHandler(cpService)
	hackathonHandler    = handler.GateHackathonHandler(SubmissionService)

	authMiddleware = middleware.NewAuthMiddleware(authService, jwtService)

	dashboards = handler.DashboardController(database)
)

func SetupRouter(r *gin.Engine) {
	api_base_url := os.Getenv("API_BASE_URL")
	router := r.Group(api_base_url)

	router.GET("/ping", authHandler.Ping)

	// Auth Router /auth
	SetupAuthRouter(router)

	mustAuth := router.Group("")
	mustAuth.Use(authMiddleware.JwtAuthMiddleware)
	mustAuth.GET("mustauth/ping", authHandler.Ping)

	// Profile Route
	{
		profile := mustAuth.Group("profile")
		profile.GET("my", userHandler.GetMyProfile)
		profile.POST("my", userHandler.UpdateMyProfile)
		profile.GET("events", userHandler.GetEvents)

		router.GET("profile/all/:start_date/:end_date/:count/:page", userHandler.GetAllUser)
	}

	mustUpdatedProfile := mustAuth.Group("")
	mustUpdatedProfile.Use(authMiddleware.MustUpdatedUserProfile)
	mustUpdatedProfile.GET("mustauth/authupdate/ping", authHandler.Ping)

	// Team Registrasion
	{
		teamRegistration := mustUpdatedProfile.Group("team/registration")
		teamRegistration.POST("hackathon", registrationHandler.RegistrationHackathonTeam)
		teamRegistration.POST("cp", registrationHandler.RegistrationCPTeam)
		teamRegistration.GET("find/:join_code", registrationHandler.FindTeam)
		teamRegistration.POST("join/:join_code", registrationHandler.UserJoinTeam)
	}

	// Profile Route
	// router.POST("profile/post", profileHandler.Create)
	// router.GET("profile/get", profileHandler.GetProfile)

	// admin_api_base_url := os.Getenv("ADMIN_API_BASE_URL")
	// admin_router := r.Group(admin_api_base_url)

	{
		newsletter := router.Group("/newsletter")

		newsletter.GET("/:id", newsletterHandler.GetNewsLetter)

		// newsletter admin
		newsletter.Use(authMiddleware.JwtAuthMiddleware)
		newsletter.Use(authMiddleware.MustAdmin)
		newsletter.POST("/", newsletterHandler.CreateNewsletter)
		newsletter.PUT("/:id", newsletterHandler.UpdateNewsLetter)
		newsletter.DELETE("/:id", newsletterHandler.DeleteNewsLetter)
	}

	{
		dashboard := router.Group("/dashboard")
		dashboard.Use(authMiddleware.JwtAuthMiddleware)
		dashboardUnauth := router.Group("/dashboard")

		dashboardUnauth.GET("/:acara/:start_date/:end_date/:count/:page", dashboards.GetAllDashboard)
		dashboard.DELETE("/:acara/:id", dashboards.Delete)
		dashboard.PUT("/:acara/:id", dashboards.Update)
		dashboard.GET("/events/:id_user", dashboards.GetEvent)
	}

	{
		submissionHandler := router.Group("/submission")
		submissionHandler.POST("/hackaton/:stage/:join_code", hackathonHandler.SubmissionHackaton)
		submissionHandler.GET("/hackaton/:join_code", hackathonHandler.HackathonStageStatus)
	}

	{
		cp := router.Group("/cp")
		cp.GET("/:join_code", cpHandler.GetDetail)
	}
}
