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
	GetAuditLogStats(c *gin.Context)
}

func NewAuditLogHandler(auditLogService service.AuditLogService) AuditLogHandler {
	return &auditLogHandler{
		auditLogService: auditLogService,
	}
}

// GetAllAuditLogs retrieves all audit logs with pagination
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

	role := c.Query("role")
	query := c.Query("q")
	offset := (page - 1) * limit

	logs, total, err := h.auditLogService.GetAllActivityLogs(limit, offset, role, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Failed to retrieve audit logs"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Successfully retrieved audit logs", gin.H{
		"logs":        logs,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}))
}

// GetUserAuditLogs retrieves audit logs for a specific user
func (h *auditLogHandler) GetUserAuditLogs(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"user_id": {"IS_INVALID"}}))
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
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Failed to retrieve audit logs"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Successfully retrieved user audit logs", gin.H{
		"logs":        logs,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}))
}

// GetAuditLogsByDateRange retrieves audit logs within a date range
func (h *auditLogHandler) GetAuditLogsByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "start_date and end_date are required"))
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Invalid start_date format. Use YYYY-MM-DD"))
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Invalid end_date format. Use YYYY-MM-DD"))
		return
	}

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

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"user_id": {"IS_INVALID"}}))
			return
		}

		logs, total, err = h.auditLogService.GetUserActivityLogsByDateRange(userID, startDate, endDate, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Failed to retrieve audit logs"))
			return
		}
	} else {
		logs, total, err = h.auditLogService.GetActivityLogsByDateRange(startDate, endDate, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Failed to retrieve audit logs"))
			return
		}
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Successfully retrieved audit logs", gin.H{
		"logs":        logs,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}))
}

func (h *auditLogHandler) GetAuditLogStats(c *gin.Context) {
	stats, err := h.auditLogService.GetAuditLogStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Failed to retrieve audit log stats"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Successfully retrieved audit log stats", stats))
}
