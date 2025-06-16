package router

import (
	"gcw/dto"
	"gcw/middleware"

	"github.com/gin-gonic/gin"
)

// /auth
func SetupAuthRouter(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	// Authentication Route
	auth.POST("validate-google-id-token", authHandler.ValidateGoogleIdToken)
	auth.POST("refresh-token", authHandler.RefreshToken)
	auth.POST("send-mail-test", authHandler.SendEmailVerificationExample)

	auth.POST("login", middleware.ValidateDTO(&dto.LoginDTO{}), authHandler.Login)
	auth.POST("registration", middleware.ValidateDTO(&dto.RegisterDTO{}), authHandler.Registration)

	mustAuth := auth.Group("")
	mustAuth.Use(authMiddleware.JwtAuthMiddleware)
	// logout route
}
