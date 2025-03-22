package helper

import (
	"os"

	"github.com/gin-gonic/gin"
)

var (
	cookieOrigin = os.Getenv("COOKIE_ORIGIN")
)

func SetTokenRefreshCookie(c *gin.Context, token string, refreshToken string) {
	c.SetCookie("gcw_api_token", token, 60*60*24*7, "/", cookieOrigin, false, true)
	c.SetCookie("gcw_api_refresh_token", refreshToken, 60*60*24*7, "/", cookieOrigin, false, true)
}

func RemoveTokenRefreshCookie(c *gin.Context) {
	c.SetCookie("gcw_api_token", "", -1, "/", cookieOrigin, false, true)
	c.SetCookie("gcw_api_refresh_token", "", -1, "/", cookieOrigin, false, true)
}
