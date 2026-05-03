package middleware

import (
	"bytes"
	"encoding/json"
	"gcw/entity"
	"gcw/service"
	"io"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type auditMiddleware struct {
	auditLogService service.AuditLogService
}

type AuditMiddlewareI interface {
	AuditLogMiddleware(*gin.Context)
}

func NewAuditMiddleware(als service.AuditLogService) AuditMiddlewareI {
	return &auditMiddleware{
		auditLogService: als,
	}
}

// AuditLogMiddleware records user activity for authenticated requests (excluding GET/HEAD)
func (m *auditMiddleware) AuditLogMiddleware(c *gin.Context) {
	// Only log authenticated requests
	userInterface, exists := c.Get("user")
	if !exists {
		c.Next()
		return
	}

	user, ok := userInterface.(*entity.User)
	if !ok {
		c.Next()
		return
	}

	// Skip logging for certain endpoints (like health checks)
	if shouldSkipAuditLog(c.Request.URL.Path) {
		c.Next()
		return
	}

	// Skip logging for authentication and registration endpoints
	if shouldSkipAuthAuditLog(c.Request.URL.Path) {
		c.Next()
		return
	}

	// Skip logging for GET and HEAD requests
	if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
		c.Next()
		return
	}

	// Capture request body
	var requestBody interface{}
	if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
		var bodyBytes []byte
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Parse request body
		if len(bodyBytes) > 0 {
			contentType := c.GetHeader("Content-Type")
			if strings.Contains(contentType, "application/json") {
				var jsonBody map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
					requestBody = jsonBody
				}
			} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
				if err := c.Request.ParseForm(); err == nil {
					formData := url.Values{}
					for key, values := range c.Request.PostForm {
						formData[key] = values
					}
					requestBody = formData
				}
			}
		}
	}

	// Get client IP address
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Create response writer wrapper to capture status code
	responseWriter := &ResponseWriter{
		ResponseWriter: c.Writer,
		statusCode:     200,
	}
	c.Writer = responseWriter

	// Process request
	c.Next()

	// Record audit log with response code
	err := m.auditLogService.RecordActivityWithResponse(
		user.ID,
		c.Request.Method,
		c.Request.URL.Path,
		requestBody,
		responseWriter.statusCode,
		nil,
		ipAddress,
		userAgent,
	)

	if err != nil {
		// Log error but don't affect the response
		_ = err
	}
}

// ResponseWriter wraps gin.ResponseWriter to capture status code
type ResponseWriter struct {
	gin.ResponseWriter
	statusCode int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// shouldSkipAuditLog returns true if the endpoint should not be logged
func shouldSkipAuditLog(path string) bool {
	skipPaths := []string{
		"/health",
		"/metrics",
		"/api/v1/health",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	return false
}

// shouldSkipAuthAuditLog returns true for auth and registration routes that should not be audited
func shouldSkipAuthAuditLog(path string) bool {
	return strings.Contains(path, "/auth/") || strings.HasSuffix(path, "/auth")
}
