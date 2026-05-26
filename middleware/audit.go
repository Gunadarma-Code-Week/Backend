package middleware

import (
	"bytes"
	"encoding/json"
	"gcw/entity"
	"gcw/service"
	"fmt"
	"io"
	"net/url"
	"strconv"
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
					sanitizeSensitiveFields(jsonBody)
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
	description, targetResourceID := m.generateDescription(c.Request.Method, c.Request.URL.Path, requestBody, user, c)
	targetEmail := c.GetString("target_email")
	targetName := c.GetString("target_name")

	err := m.auditLogService.RecordActivityWithResponse(
		user.ID,
		user.Email,
		targetEmail,
		targetName,
		targetResourceID,
		c.Request.Method,
		c.Request.URL.Path,
		description,
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
	return strings.Contains(path, "/auth/") || strings.HasSuffix(path, "/auth") || strings.Contains(path, "/profile/change-password")
}

// sensitiveFields is the list of JSON keys whose values will be redacted in audit logs
var sensitiveFields = []string{
	"password",
	"domjudge_password",
	"token",
	"access_token",
	"refresh_token",
	"secret",
}

// sanitizeSensitiveFields replaces sensitive field values with "[REDACTED]" in-place
func sanitizeSensitiveFields(body map[string]interface{}) {
	for _, field := range sensitiveFields {
		if _, ok := body[field]; ok {
			body[field] = "[REDACTED]"
		}
	}
}

// generateDescription generates a human-readable Indonesian description and extracts target resource ID
func (m *auditMiddleware) generateDescription(method, path string, body interface{}, user *entity.User, c *gin.Context) (string, uint64) {
	var targetResourceID uint64
	
	// Check context for target_id set by handler (common for registration/post actions)
	if cid, exists := c.Get("target_id"); exists {
		if id, ok := cid.(uint64); ok {
			targetResourceID = id
		} else if id, ok := cid.(int64); ok {
			targetResourceID = uint64(id)
		} else if id, ok := cid.(int); ok {
			targetResourceID = uint64(id)
		}
	}
	// User identifier: use email if available, otherwise "User"
	userLabel := "User"
	if user != nil && user.Email != "" {
		userLabel = user.Email
	}

	// Target name/email from context if set by handler
	targetName := c.GetString("target_name")
	targetEmail := c.GetString("target_email")

	// Normalize path segments
	segments := strings.Split(strings.Trim(path, "/"), "/")

	// Dashboard actions: /api/v1/gcw/resources/dashboard/{acara}/{id}
	if len(segments) >= 2 {
		n := len(segments)
		acara := segments[n-2]
		id := segments[n-1]
		if acara == "hackaton" || acara == "hackathon" || acara == "cp" || acara == "ctf" || acara == "seminar" {
			acaraLabel := map[string]string{
				"hackaton":  "Hackathon",
				"hackathon": "Hackathon",
				"cp":        "Competitive Programming",
				"ctf":       "CTF",
				"seminar":   "Seminar",
			}[acara]
			targetLabel := "dengan ID " + id
			targetResourceID, _ = strconv.ParseUint(id, 10, 64)
			if targetName != "" {
				targetLabel = "dengan nama " + targetName
			}

			switch method {
			case "DELETE":
				description := userLabel + " menghapus tim " + acaraLabel + " " + targetLabel
				if bodyMap, ok := body.(map[string]interface{}); ok {
					if alasan, hasAlasan := bodyMap["alasan"]; hasAlasan {
						if alasanStr, ok := alasan.(string); ok && alasanStr != "" {
							description += " dengan alasan: " + alasanStr
						}
					}
				}
				return description, targetResourceID
			case "PUT", "PATCH":
				// Detect specific field changes from request body
				if bodyMap, ok := body.(map[string]interface{}); ok {
					var changes []string
					if _, hasPass := bodyMap["password"]; hasPass {
						changes = append(changes, "password DomJudge")
					}
					if _, hasUser := bodyMap["username"]; hasUser {
						changes = append(changes, "username DomJudge")
					}
					if val, hasStage := bodyMap["stage"]; hasStage {
						if stageStr, ok := val.(string); ok {
							changes = append(changes, fmt.Sprintf("stage ke \"%s\"", stageStr))
						} else {
							changes = append(changes, "stage")
						}
					}
					if val, hasNama := bodyMap["nama_tim"]; hasNama {
						if namaStr, ok := val.(string); ok {
							changes = append(changes, fmt.Sprintf("nama tim ke \"%s\"", namaStr))
						} else {
							changes = append(changes, "nama tim")
						}
					}
					if len(changes) > 0 {
						return userLabel + " memperbarui " + strings.Join(changes, ", ") + " tim " + acaraLabel + " " + targetLabel, targetResourceID
					}
				}
				return userLabel + " memperbarui data tim " + acaraLabel + " " + targetLabel, targetResourceID
			}
		}
	}

	// Admin user management: /admin/users/{id}
	if strings.Contains(path, "/admin/users/") {
		id := segments[len(segments)-1]
		targetLabel := "dengan ID " + id
		targetResourceID, _ = strconv.ParseUint(id, 10, 64)
		if targetEmail != "" {
			targetLabel = "dengan email " + targetEmail
		} else if targetName != "" {
			targetLabel = "dengan nama " + targetName
		}
		switch method {
		case "PUT", "PATCH":
			if bodyMap, ok := body.(map[string]interface{}); ok {
				var changes []string
				if val, hasVerified := bodyMap["data_has_verified"]; hasVerified {
					if verified, ok := val.(bool); ok {
						if verified {
							changes = append(changes, "memverifikasi data")
						} else {
							changes = append(changes, "membatalkan verifikasi data")
						}
					}
				}
				if val, hasProfileUpdated := bodyMap["profile_has_updated"]; hasProfileUpdated {
					if updated, ok := val.(bool); ok {
						if updated {
							changes = append(changes, "menandai profil lengkap")
						} else {
							changes = append(changes, "menandai profil belum lengkap")
						}
					}
				}
				
				if len(changes) > 0 {
					return userLabel + " " + strings.Join(changes, " dan ") + " user " + targetLabel, targetResourceID
				}
			}
			return userLabel + " memperbarui data user " + targetLabel, targetResourceID
		case "DELETE":
			description := userLabel + " menghapus user " + targetLabel
			if bodyMap, ok := body.(map[string]interface{}); ok {
				if alasan, hasAlasan := bodyMap["alasan"]; hasAlasan {
					if alasanStr, ok := alasan.(string); ok && alasanStr != "" {
						description += " dengan alasan: " + alasanStr
					}
				}
			}
			return description, targetResourceID
		}
	}

	// Team registration
	if strings.Contains(path, "/team/registration/") {
		switch method {
		case "POST":
			if strings.Contains(path, "/join/") {
				if targetName != "" {
					return userLabel + " bergabung ke tim: " + targetName, targetResourceID
				}
				joinCode := segments[len(segments)-1]
				return userLabel + " bergabung ke tim dengan join code: " + joinCode, targetResourceID
			}
			
			// For new registrations, try to get the team name from targetName or body
			actualTeamName := targetName
			if actualTeamName == "" {
				if bodyMap, ok := body.(map[string]interface{}); ok {
					if teamName, hasName := bodyMap["team_name"]; hasName {
						if nameStr, ok := teamName.(string); ok && nameStr != "" {
							actualTeamName = nameStr
						}
					}
				}
			}

			if actualTeamName != "" {
				eventLabel := "tim baru"
				if strings.Contains(path, "/hackathon") {
					eventLabel = "tim Hackathon baru"
				} else if strings.Contains(path, "/cp") {
					eventLabel = "tim Competitive Programming baru"
				} else if strings.Contains(path, "/ctf") {
					eventLabel = "tim CTF baru"
				}
				return userLabel + " mendaftarkan " + eventLabel + ": " + actualTeamName, targetResourceID
			}
			return userLabel + " mendaftarkan tim baru", targetResourceID
		}
	}

	// Profile update
	if strings.Contains(path, "/profile/my") && method == "POST" {
		description := userLabel + " memperbarui profil"
		auditChanges := c.GetString("audit_changes")
		if auditChanges != "" {
			description += " (" + auditChanges + ")"
		}
		return description, targetResourceID
	}

	// Seminar
	if strings.Contains(path, "/seminar/join") && method == "POST" {
		return userLabel + " mendaftar seminar", targetResourceID
	}
	if strings.Contains(path, "/seminar/admin/add-participant") && method == "POST" {
		return userLabel + " menambahkan peserta seminar", targetResourceID
	}

	// Submission
	if strings.Contains(path, "/submission/") && method == "POST" {
		stage := "submission"
		// Path format: /api/v1/gcw/resources/submission/hackaton/:stage/:join_code
		for i, segment := range segments {
			if segment == "submission" && i+2 < len(segments) {
				stage = segments[i+2]
				break
			}
		}

		stageLabel := "submission"
		switch stage {
		case "stage1":
			stageLabel = "proposal (Stage 1)"
		case "stage2":
			stageLabel = "pitch deck (Stage 2)"
		case "final":
			stageLabel = "video / presentasi final"
		}

		desc := userLabel + " mengirimkan " + stageLabel
		if targetName != "" {
			desc += " untuk tim: " + targetName
		}
		return desc, targetResourceID
	}

	// Payment details update
	if strings.Contains(path, "/payment/update-team-details") && method == "POST" {
		auditChanges := c.GetString("audit_changes")
		actualTeamName := c.GetString("target_name")
		if actualTeamName == "" {
			if bodyMap, ok := body.(map[string]interface{}); ok {
				if teamName, hasName := bodyMap["team_name"]; hasName {
					if nameStr, ok := teamName.(string); ok && nameStr != "" {
						actualTeamName = nameStr
					}
				}
			}
		}

		desc := userLabel + " memperbarui detail data tim"
		if actualTeamName != "" {
			desc += ": " + actualTeamName
		}
		if auditChanges != "" {
			desc += " (Perubahan: " + auditChanges + ")"
		}
		return desc, targetResourceID
	}

	// System Settings Update
	if strings.Contains(path, "/settings") && method == "PUT" {
		auditChanges := c.GetString("audit_changes")
		desc := userLabel + " memperbarui konfigurasi pengumpulan berkas"
		if auditChanges != "" {
			desc += " (" + auditChanges + ")"
		}
		return desc, 0
	}

	// Timeline Update
	if strings.Contains(path, "/admin/timeline") {
		action := "mengakses"
		switch method {
		case "POST":
			action = "membuat"
		case "PUT", "PATCH":
			action = "memperbarui"
		case "DELETE":
			action = "menghapus"
		}
		
		desc := userLabel + " " + action + " data timeline"
		if len(segments) > 0 && segments[len(segments)-1] != "timeline" {
			id := segments[len(segments)-1]
			targetResourceID, _ = strconv.ParseUint(id, 10, 64)
			desc += " dengan ID " + id
		}
		return desc, targetResourceID
	}

	// Fallback: generic description
	actionMap := map[string]string{
		"POST":   "membuat",
		"PUT":    "memperbarui",
		"PATCH":  "memperbarui",
		"DELETE": "menghapus",
	}
	action, ok := actionMap[method]
	if !ok {
		action = "mengakses"
	}
	desc := userLabel + " " + action + " resource pada " + path
	if targetEmail != "" {
		desc += " dengan email " + targetEmail
	}
	return desc, targetResourceID
}
