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
	MustAdmin(*gin.Context)
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

	// Safely extract user ID from claims
	idVal, ok := claims["id"]
	if !ok || idVal == nil {
		c.JSON(401, helper.CreateErrorResponse("error", "invalid token: missing id"))
		c.Abort()
		return
	}

	var idUser uint64
	switch v := idVal.(type) {
	case float64:
		idUser = uint64(v)
	case int:
		idUser = uint64(v)
	case int64:
		idUser = uint64(v)
	case uint64:
		idUser = v
	default:
		c.JSON(401, helper.CreateErrorResponse("error", "invalid token: id is not a number"))
		c.Abort()
		return
	}

	if idUser == 0 {
		c.JSON(401, helper.CreateErrorResponse("error", "invalid token: id is zero"))
		c.Abort()
		return
	}

	user, err := m.authService.GetUserById(idUser)
	if err != nil {
		// Check if user's account has been soft deleted
		isDeleted, reason := m.authService.IsUserSoftDeleted(idUser)
		if isDeleted {
			msg := "Akun Anda telah dihapus oleh admin"
			if reason != "" {
				msg += " dengan alasan: " + reason
			}

			logging.High("AuthMiddleware.JwtAuthMiddleware", "USER_DELETED",
				"User account has been soft-deleted, forcing logout")
			c.JSON(401, gin.H{
				"status":  "error",
				"message": msg,
				"code":    "USER_DELETED",
			})
			c.Abort()
			return
		}

		logging.High("AuthMiddleware.JwtAuthMiddleware", "USER_NOT_FOUND", err.Error())
		c.JSON(401, helper.CreateErrorResponse("error", "user not found"))
		c.Abort()
		return
	}

	// Check if user's team has been soft deleted
	if user.IDTeam != nil {
		isDeleted, teamEvent := m.authService.IsUserTeamDeleted(user)
		if isDeleted {
			logging.High("AuthMiddleware.JwtAuthMiddleware", "TEAM_DELETED",
				"User team has been deleted, forcing logout")
			c.JSON(401, gin.H{
				"status":  "error",
				"message": "Tim Anda telah dihapus oleh admin",
				"code":    "TEAM_DELETED",
				"event":   teamEvent,
			})
			c.Abort()
			return
		}
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

func (m *authMiddleware) MustAdmin(c *gin.Context) {
	userAuth, ok := c.MustGet("user").(*entity.User)

	if !ok || userAuth.Role != "admin" {
		c.JSON(403, helper.CreateErrorResponse("error", "access denied"))
		c.Abort()
		return
	}

	c.Next()
}
