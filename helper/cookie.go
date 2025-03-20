package helper

import (
	"os"

	"github.com/gin-gonic/gin"
)

func SetTokenRefreshCookie(c *gin.Context, token string, refreshToken string) {
	origin := os.Getenv("COOKIE_ORIGIN")

	c.SetCookie("gcw_api_token", token, 60*60*24*7, "/", origin, false, true)
	c.SetCookie("gcw_api_refresh_token", refreshToken, 60*60*24*7, "/", origin, false, true)
}
