package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimiter returns a gin.HandlerFunc for rate limiting.
func RateLimiter() gin.HandlerFunc {
	// Define the rate (100 requests per minute)
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}

	// Create an in-memory store
	store := memory.NewStore()

	// Create the limiter instance
	instance := limiter.New(store, rate)

	// Create the middleware
	middleware := mgin.NewMiddleware(instance, mgin.WithErrorHandler(func(c *gin.Context, err error) {
		log.Printf("Rate limit error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error during rate limiting",
		})
	}))

	return middleware
}
