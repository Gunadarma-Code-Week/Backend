package service

import (
	"encoding/json"
	"gcw/entity"
	"gcw/repository"
	"time"
)

type auditLogService struct {
	auditLogRepository repository.AuditLogRepository
}

type AuditLogService interface {
	RecordActivity(userID uint64, method string, endpoint string, requestBody interface{}, ipAddress string, userAgent string) error
	RecordActivityWithResponse(userID uint64, method string, endpoint string, requestBody interface{}, responseCode int, responseBody interface{}, ipAddress string, userAgent string) error
	RecordActivityWithError(userID uint64, method string, endpoint string, requestBody interface{}, responseCode int, errorMessage string, ipAddress string, userAgent string) error
	GetUserActivityLogs(userID uint64, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetEndpointActivityLogs(endpoint string, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetActivityLogsByDateRange(startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetUserActivityLogsByDateRange(userID uint64, startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetAllActivityLogs(limit int, offset int) ([]*entity.AuditLog, int64, error)
}

func NewAuditLogService(repo repository.AuditLogRepository) AuditLogService {
	return &auditLogService{
		auditLogRepository: repo,
	}
}

// RecordActivity records user activity without response data
func (s *auditLogService) RecordActivity(userID uint64, method string, endpoint string, requestBody interface{}, ipAddress string, userAgent string) error {
	var requestJSON json.RawMessage
	if requestBody != nil {
		data, err := json.Marshal(requestBody)
		if err != nil {
			requestJSON = []byte("{}")
		} else {
			requestJSON = data
		}
	} else {
		requestJSON = []byte("{}")
	}

	auditLog := &entity.AuditLog{
		UserID:      userID,
		Method:      method,
		Endpoint:    endpoint,
		RequestBody: requestJSON,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}

	return s.auditLogRepository.Create(auditLog)
}

// RecordActivityWithResponse records user activity with response data
func (s *auditLogService) RecordActivityWithResponse(userID uint64, method string, endpoint string, requestBody interface{}, responseCode int, responseBody interface{}, ipAddress string, userAgent string) error {
	var requestJSON json.RawMessage
	if requestBody != nil {
		data, err := json.Marshal(requestBody)
		if err != nil {
			requestJSON = []byte("{}")
		} else {
			requestJSON = data
		}
	} else {
		requestJSON = []byte("{}")
	}

	var responseJSON json.RawMessage
	if responseBody != nil {
		data, err := json.Marshal(responseBody)
		if err != nil {
			responseJSON = []byte("{}")
		} else {
			responseJSON = data
		}
	} else {
		responseJSON = []byte("{}")
	}

	auditLog := &entity.AuditLog{
		UserID:       userID,
		Method:       method,
		Endpoint:     endpoint,
		RequestBody:  requestJSON,
		ResponseCode: responseCode,
		ResponseBody: responseJSON,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}

	return s.auditLogRepository.Create(auditLog)
}

// RecordActivityWithError records user activity with error information
func (s *auditLogService) RecordActivityWithError(userID uint64, method string, endpoint string, requestBody interface{}, responseCode int, errorMessage string, ipAddress string, userAgent string) error {
	var requestJSON json.RawMessage
	if requestBody != nil {
		data, err := json.Marshal(requestBody)
		if err != nil {
			requestJSON = []byte("{}")
		} else {
			requestJSON = data
		}
	} else {
		requestJSON = []byte("{}")
	}

	auditLog := &entity.AuditLog{
		UserID:       userID,
		Method:       method,
		Endpoint:     endpoint,
		RequestBody:  requestJSON,
		ResponseCode: responseCode,
		ErrorMessage: &errorMessage,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}

	return s.auditLogRepository.Create(auditLog)
}

// GetUserActivityLogs retrieves activity logs for a specific user
func (s *auditLogService) GetUserActivityLogs(userID uint64, limit int, offset int) ([]*entity.AuditLog, int64, error) {
	return s.auditLogRepository.FindByUserID(userID, limit, offset)
}

// GetEndpointActivityLogs retrieves activity logs for a specific endpoint
func (s *auditLogService) GetEndpointActivityLogs(endpoint string, limit int, offset int) ([]*entity.AuditLog, int64, error) {
	return s.auditLogRepository.FindByEndpoint(endpoint, limit, offset)
}

// GetActivityLogsByDateRange retrieves activity logs within a date range
func (s *auditLogService) GetActivityLogsByDateRange(startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error) {
	return s.auditLogRepository.FindByDateRange(startDate, endDate, limit, offset)
}

// GetUserActivityLogsByDateRange retrieves activity logs for a user within a date range
func (s *auditLogService) GetUserActivityLogsByDateRange(userID uint64, startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error) {
	return s.auditLogRepository.FindByUserIDAndDateRange(userID, startDate, endDate, limit, offset)
}

// GetAllActivityLogs retrieves all activity logs
func (s *auditLogService) GetAllActivityLogs(limit int, offset int) ([]*entity.AuditLog, int64, error) {
	return s.auditLogRepository.FindAll(limit, offset)
}
