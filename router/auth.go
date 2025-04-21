package router

import "github.com/gin-gonic/gin"

// /auth
func SetupAuthRouter(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	// Authentication Route
	auth.POST("validate-google-id-token", authHandler.ValidateGoogleIdToken)
	auth.POST("refresh-token", authHandler.RefreshToken)
	auth.POST("send-mail-test", authHandler.SendEmailVerificationExample)

	mustAuth := auth.Group("")
	mustAuth.Use(authMiddleware.JwtAuthMiddleware)
	// logout route
}
