package handler

import (
	"gcw/helper"
	"gcw/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type auditLogHandler struct {
	auditLogService service.AuditLogService
}

type AuditLogHandler interface {
	GetAllAuditLogs(c *gin.Context)
	GetUserAuditLogs(c *gin.Context)
	GetAuditLogsByDateRange(c *gin.Context)
}

func NewAuditLogHandler(auditLogService service.AuditLogService) AuditLogHandler {
	return &auditLogHandler{
		auditLogService: auditLogService,
	}
}

// GetAllAuditLogs retrieves all audit logs with pagination
// @Summary Get all audit logs
// @Description Get all audit logs with pagination. Requires admin role.
// @Tags AuditLog
// @Security Bearer
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/audit-logs [get]
func (h *auditLogHandler) GetAllAuditLogs(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := (page - 1) * limit

	logs, total, err := h.auditLogService.GetAllActivityLogs(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", "Failed to retrieve audit logs"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Successfully retrieved audit logs", gin.H{
		"logs":       logs,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}))
}

// GetUserAuditLogs retrieves audit logs for a specific user
// @Summary Get user audit logs
// @Description Get audit logs for a specific user with pagination.
// @Tags AuditLog
// @Security Bearer
// @Param user_id path int true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/audit-logs/user/{user_id} [get]
func (h *auditLogHandler) GetUserAuditLogs(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "Invalid user ID"))
		return
	}

	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := (page - 1) * limit

	logs, total, err := h.auditLogService.GetUserActivityLogs(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", "Failed to retrieve audit logs"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Successfully retrieved user audit logs", gin.H{
		"logs":       logs,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}))
}

// GetAuditLogsByDateRange retrieves audit logs within a date range
// @Summary Get audit logs by date range
// @Description Get audit logs within a specific date range with optional user filter.
// @Tags AuditLog
// @Security Bearer
// @Param start_date query string true "Start date (YYYY-MM-DD format)"
// @Param end_date query string true "End date (YYYY-MM-DD format)"
// @Param user_id query int false "Filter by user ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Router /admin/audit-logs/date-range [get]
func (h *auditLogHandler) GetAuditLogsByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "start_date and end_date are required"))
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "Invalid start_date format. Use YYYY-MM-DD"))
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "Invalid end_date format. Use YYYY-MM-DD"))
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(time.Hour * 24)

	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := (page - 1) * limit

	var logs interface{}
	var total int64

	// Check if user_id is provided
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "Invalid user ID"))
			return
		}

		logs, total, err = h.auditLogService.GetUserActivityLogsByDateRange(userID, startDate, endDate, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", "Failed to retrieve audit logs"))
			return
		}
	} else {
		logs, total, err = h.auditLogService.GetActivityLogsByDateRange(startDate, endDate, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", "Failed to retrieve audit logs"))
			return
		}
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Successfully retrieved audit logs", gin.H{
		"logs":       logs,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}))
}
