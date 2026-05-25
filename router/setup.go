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
	auditLogRepository     = repository.NewAuditLogRepository(database)
	timelineRepository     = repository.NewTimelineRepository(database)

	jwtService          = service.NewJwtService()
	emailService        = service.NewEmailService()
	domJudgeService     = service.NewDomJudgeService()
	stellarService      = service.NewStellarService()
	authService         = service.NewAuthService(userRepository, database)
	userService         = service.NewUserService(userRepository)
	registrationService = service.NewRegistrationService(
		registrationRepository,
		domJudgeService,
	)
	newsletterService = service.NewNewsletterService(newsletterRepository)
	SubmissionService = service.NewSubmissionService(database)
	cpService         = service.NewCpService(database)
	ctfService        = service.NewCtfService(database)
	seminarService    = service.NewSeminarService(database)
	auditLogService   = service.NewAuditLogService(auditLogRepository, userRepository, stellarService)
	timelineService   = service.NewTimelineService(timelineRepository)

	authHandler         = handler.NewAuthHandler(authService, jwtService, emailService)
	userHandler         = handler.NewUserHandler(userService)
	registrationHandler = handler.GateRegistrationHandler(registrationService, userService)
	paymentHandler      = handler.NewPaymentHandler(registrationService)
	// newsletterHandler   = handler.NewNewsletterHandler(newsletterService)
	submissionHandler = handler.GateHackathonHandler(SubmissionService)
	cpHandler         = handler.GateCompetitiveHandler(cpService)
	ctfHandler        = handler.NewCTFHandler(ctfService)
	hackathonHandler  = handler.GateHackathonHandler(SubmissionService)
	seminarHandler    = handler.NewSeminarHandler(seminarService)
	auditLogHandler   = handler.NewAuditLogHandler(auditLogService)
	settingService    = service.NewSystemSettingService(database)
	settingHandler    = handler.NewSystemSettingHandler(settingService)
	timelineHandler   = handler.NewTimelineHandler(timelineService)

	authMiddleware  = middleware.NewAuthMiddleware(authService, jwtService)
	auditMiddleware = middleware.NewAuditMiddleware(auditLogService)

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
	mustAuth.Use(auditMiddleware.AuditLogMiddleware)
	mustAuth.GET("mustauth/ping", authHandler.Ping)

	// Profile Route
	{
		profile := mustAuth.Group("profile")
		profile.GET("my", userHandler.GetMyProfile)
		profile.POST("my", userHandler.UpdateMyProfile)
		profile.POST("change-password", userHandler.ChangePassword)
		profile.GET("events", userHandler.GetEvents)

		router.GET("profile/all/:start_date/:end_date/:count/:page", userHandler.GetAllUser)
	}

	mustUpdatedProfile := mustAuth.Group("")
	mustUpdatedProfile.Use(authMiddleware.MustUpdatedUserProfile)
	mustUpdatedProfile.GET("mustauth/authupdate/ping", authHandler.Ping)

	mustVerified := mustUpdatedProfile.Group("")
	mustVerified.Use(authMiddleware.MustVerifiedUser)

	// Team Registrasion
	{
		teamRegistration := mustVerified.Group("team/registration")
		teamRegistration.POST("hackathon", registrationHandler.RegistrationHackathonTeam)
		teamRegistration.POST("cp", registrationHandler.RegistrationCPTeam)
		teamRegistration.POST("ctf", registrationHandler.RegistrationCTFTeam)
		teamRegistration.GET("find/:join_code", registrationHandler.FindTeam)
		teamRegistration.POST("join/:join_code", registrationHandler.UserJoinTeam)
	}



	// Authenticated Payment Actions
	{
		authPayment := mustAuth.Group("/payment")
		authPayment.POST("/update-team-details", paymentHandler.UpdateTeamDetails)
	}

	// Profile Route
	// router.POST("profile/post", profileHandler.Create)
	// router.GET("profile/get", profileHandler.GetProfile)

	// admin_api_base_url := os.Getenv("ADMIN_API_BASE_URL")
	// admin_router := r.Group(admin_api_base_url)

	// {
	// 	newsletter := router.Group("/newsletter")

	// 	newsletter.GET("/:id", newsletterHandler.GetNewsLetter)

	// 	// newsletter admin
	// 	newsletter.Use(authMiddleware.JwtAuthMiddleware)
	// 	newsletter.Use(authMiddleware.MustAdmin)
	// 	newsletter.POST("/", newsletterHandler.CreateNewsletter)
	// 	newsletter.PUT("/:id", newsletterHandler.UpdateNewsLetter)
	// 	newsletter.DELETE("/:id", newsletterHandler.DeleteNewsLetter)
	// }

	{
		dashboard := mustAuth.Group("/dashboard")
		dashboard.Use(authMiddleware.MustAdmin)
		dashboard.GET("/:acara/:count/:page", dashboards.GetAllDashboard)
		dashboard.DELETE("/:acara/:id", dashboards.Delete)
		dashboard.PUT("/:acara/:id", dashboards.Update)
	}

	// Admin User Management Routes
	{
		adminUsers := mustAuth.Group("/admin/users")
		adminUsers.Use(authMiddleware.MustAdmin)

		adminUsers.GET("", userHandler.AdminGetAllUsers)
		adminUsers.GET("/:id", userHandler.AdminGetUserById)
		adminUsers.PUT("/:id", userHandler.AdminUpdateUser)
		adminUsers.DELETE("/:id", userHandler.AdminDeleteUser)
		adminUsers.GET("/analytics/growth", userHandler.AdminGetUserGrowthAnalytics)
	}

	// Admin Bulk Email and Team Names Routes
	{
		adminBulk := mustAuth.Group("/admin")
		adminBulk.Use(authMiddleware.MustAdmin)

		adminBulk.GET("/teams/names", dashboards.GetTeamNames)
		adminBulk.GET("/emails/templates", dashboards.GetEmailTemplates)
	}


	// Audit Log Routes
	{
		auditLogs := mustAuth.Group("/admin/audit-logs")
		auditLogs.Use(authMiddleware.MustAdmin)

		auditLogs.GET("", auditLogHandler.GetAllAuditLogs)
		auditLogs.GET("/stats", auditLogHandler.GetAuditLogStats)
		auditLogs.GET("/user/:user_id", auditLogHandler.GetUserAuditLogs)
		auditLogs.GET("/date-range", auditLogHandler.GetAuditLogsByDateRange)
	}

	{
		submissionHandler := mustAuth.Group("/submission")
		submissionHandler.POST("/hackaton/:stage/:join_code", hackathonHandler.SubmissionHackaton)
		submissionHandler.GET("/hackaton/:join_code", hackathonHandler.HackathonStageStatus)
	}

	{
		router.GET("/settings", settingHandler.GetSettings)

		adminSettings := mustAuth.Group("/admin/settings")
		adminSettings.Use(authMiddleware.MustAdmin)
		adminSettings.PUT("", settingHandler.UpdateSettings)
	}

	{
		router.GET("/timeline/:category", timelineHandler.GetTimelinesByCategory)

		adminTimeline := mustAuth.Group("/admin/timeline")
		adminTimeline.Use(authMiddleware.MustAdmin)
		adminTimeline.POST("", timelineHandler.CreateTimeline)
		adminTimeline.PUT("/:id", timelineHandler.UpdateTimeline)
		adminTimeline.DELETE("/:id", timelineHandler.DeleteTimeline)
	}

	{
		cp := router.Group("/cp")
		cp.GET("/:join_code", cpHandler.GetDetail)
	}

	{
		ctf := router.Group("/ctf")
		ctf.GET("/:join_code", ctfHandler.GetDetail)
	}

	{
		seminar := mustAuth.Group("/seminar")
		// seminar.Use(authMiddleware.MustUpdatedUserProfile)
		seminar.POST("/join", seminarHandler.JoinSeminar)
		seminar.GET("/my-ticket", seminarHandler.GetMyTicket)

		// Admin route untuk melihat tiket berdasarkan ID dan menambahkan participant
		seminarAdmin := seminar.Group("")
		seminarAdmin.Use(authMiddleware.MustAdmin)
		seminarAdmin.GET("/ticket/:ticket_id", seminarHandler.GetTicketByID)
		seminarAdmin.POST("/admin/add-participant", seminarHandler.AdminAddParticipant)
	}
}
