package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func CORSMiddleware(c *gin.Context) {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: strings.Split(os.Getenv("CORS_ORIGINS"), ","),
		// AllowedMethods: []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
		// AllowedHeaders: []string{"Content-Type", "X-CSRF-Token"},
		Debug: true,
	})
	corsMiddleware.HandlerFunc(c.Writer, c.Request)
	c.Next()
}
