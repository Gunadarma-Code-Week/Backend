package middleware

import (
	"gcw/entity"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type authMiddleware struct {
	jwtService  *service.JwtService
	authService *service.AuthService
}

type AuthMiddleware interface {
	JwtAuthMiddleware(*gin.Context)
	MustUpdatedUserProfile(*gin.Context)
}

func NewAuthMiddleware(as *service.AuthService, js *service.JwtService) AuthMiddleware {
	return &authMiddleware{
		authService: as,
		jwtService:  js,
	}
}

func (m *authMiddleware) JwtAuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(401, helper.CreateErrorResponse("error", "token is required"))
		c.Abort()
		return
	}

	// replace Bearer with empty string
	token = strings.Replace(token, "Bearer ", "", 1)

	claims, err := m.jwtService.GetClaimsByToken(token)
	if err != nil {
		logging.High("AuthMiddleware.JwtAuthMiddleware", "INVALID_TOKEN", err.Error())
		c.JSON(401, helper.CreateErrorResponse("error", "invalid token"))
		c.Abort()
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(401, helper.CreateErrorResponse("error", "email not found in token"))
		c.Abort()
		return
	}

	user, err := m.authService.FindByEmail(email)
	if err != nil {
		logging.High("AuthMiddleware.JwtAuthMiddleware", "USER_NOT_FOUND", err.Error())
		c.JSON(401, helper.CreateErrorResponse("error", "user not found"))
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}

func (m *authMiddleware) MustUpdatedUserProfile(c *gin.Context) {
	userAuth := c.MustGet("user").(*entity.User)

	if !userAuth.ProfileHasUpdated {
		c.JSON(400, helper.CreateErrorResponse("error", "profile has not been updated"))
		c.Abort()
		return
	}

	c.Next()
}
