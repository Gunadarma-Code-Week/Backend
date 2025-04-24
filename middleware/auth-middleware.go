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
	// get token from header
	token := c.Request.Header.Get("Authorization")
	if len(token) < 7 {
		c.JSON(401, helper.CreateErrorResponse("error", "token is required"))
		c.Abort()
		return
	}

	token = token[7:]
	if token == "" {
		c.JSON(401, helper.CreateErrorResponse("error", "token is required"))
		c.Abort()
		return
	}

	claims, err := m.jwtService.GetClaimsByToken(token)
	if err != nil {
		logging.High("AuthMiddleware.JwtAuthMiddleware", "INVALID TOKEN", err.Error())
		c.JSON(401, helper.CreateErrorResponse("error", "invalid token"))
		c.Abort()
		return
	}

	idUser := uint64(claims["id"].(float64))
	if idUser == 0 {
		c.JSON(401, helper.CreateErrorResponse("error", "invalid token"))
		c.Abort()
		return
	}

	user, err := m.authService.GetUserById(idUser)
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
	userAuth, ok := c.MustGet("user").(*entity.User)

	if !ok || !userAuth.ProfileHasUpdated {
		c.JSON(400, helper.CreateErrorResponse("error", "profile has not been updated"))
		c.Abort()
		return
	}

	c.Next()
}
