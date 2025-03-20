package middleware

import (
	"gcw/entity"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"

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
	token, err := c.Cookie("gcw_api_token")
	if err != nil || token == "" {
		c.JSON(401, helper.CreateErrorResponse("error", "token is required"))
		c.Abort()
		return
	}

	claims, err := m.jwtService.GetClaimsByToken(token)
	if err != nil {
		refreshToken, err := c.Cookie("gcw_api_refresh_token")
		if err != nil || refreshToken == "" {
			c.JSON(401, helper.CreateErrorResponse("error", "refresh token is required"))
			c.Abort()
			return
		}

		claims, err = m.jwtService.GetClaimsByRefreshToken(refreshToken)
		if err != nil {
			c.JSON(401, helper.CreateErrorResponse("error", "invalid refresh token"))
			c.Abort()
			return
		}
		email, ok := claims["email"].(string)
		if !ok {
			c.JSON(401, helper.CreateErrorResponse("error", "email not found in refresh token"))
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

		token = m.jwtService.GenerateToken(user)
		refreshToken = m.jwtService.GenerateRefreshToken(user)

		helper.SetTokenRefreshCookie(c, token, refreshToken)

		c.Set("user", user)
		c.Next()
		return
	} else {
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
		return
	}
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
