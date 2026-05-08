package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gcw/entity"
	"gcw/repository"
	"os"
	"time"
)

type auditLogService struct {
	auditLogRepository repository.AuditLogRepository
	stellarService     StellarService
}

type AuditLogService interface {
	RecordActivity(userID uint64, method string, endpoint string, description string, requestBody interface{}, ipAddress string, userAgent string) error
	RecordActivityWithResponse(userID uint64, method string, endpoint string, description string, requestBody interface{}, responseCode int, responseBody interface{}, ipAddress string, userAgent string) error
	RecordActivityWithError(userID uint64, method string, endpoint string, description string, requestBody interface{}, responseCode int, errorMessage string, ipAddress string, userAgent string) error
	GetUserActivityLogs(userID uint64, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetEndpointActivityLogs(endpoint string, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetActivityLogsByDateRange(startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetUserActivityLogsByDateRange(userID uint64, startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetAllActivityLogs(limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetAuditLogStats() (map[string]interface{}, error)
}

func NewAuditLogService(repo repository.AuditLogRepository, stellar StellarService) AuditLogService {
	return &auditLogService{
		auditLogRepository: repo,
		stellarService:     stellar,
	}
}

// calculateHash returns SHA256 hash of a string
func (s *auditLogService) calculateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// createAuditLogWithBlockchain adds blockchain hashes and saves the audit log
func (s *auditLogService) createAuditLogWithBlockchain(auditLog *entity.AuditLog) error {
	// Ensure CreatedAt is set before hashing for consistency
	if auditLog.CreatedAt.IsZero() {
		auditLog.CreatedAt = time.Now()
	}

	// Create a dictionary for hashing that includes metadata
	metadata := map[string]string{
		"timestamp":      auditLog.CreatedAt.Format(time.RFC3339),
		"deskripsi":      auditLog.Description,
		"ip":             auditLog.IPAddress,
		"endpoint":       auditLog.Endpoint,
		"status_code":    fmt.Sprintf("%d", auditLog.ResponseCode),
		"wallet_address": s.stellarService.GetPublicKey(),
	}
	metadataJSON, _ := json.Marshal(metadata)

	// 1. Calculate hash of metadata dictionary (using JSON)
	descriptionHash := s.calculateHash(string(metadataJSON))
	auditLog.DescriptionHash = descriptionHash

	// 2. Get last audit log to link to
	lastLog, err := s.auditLogRepository.GetLastAuditLog()
	if err != nil {
		return err
	}

	lastBlockchainHash := ""
	if lastLog != nil {
		lastBlockchainHash = lastLog.BlockchainHash
	}

	// 3. Calculate current blockchain hash: hash(descriptionHash + lastBlockchainHash)
	// This creates a chain of integrity
	currentBlockchainHash := s.calculateHash(descriptionHash + lastBlockchainHash)
	auditLog.BlockchainHash = currentBlockchainHash

	err = s.auditLogRepository.Create(auditLog)
	if err != nil {
		return err
	}

	// 4. Submit to Stellar blockchain in background
	go func(id uint64, hash string) {
		txHash, authorAddress, err := s.stellarService.SendAuditHash(hash)
		if err != nil {
			fmt.Printf("Stellar background task error: %v\n", err)
			// Even if it fails, we might still have the author address if the secret was valid
			if authorAddress != "" {
				db := s.auditLogRepository.GetDB()
				db.Model(&entity.AuditLog{}).Where("id = ?", id).Update("author_address", authorAddress)
			}
			return
		}

		// Update database with transaction hash and author address
		db := s.auditLogRepository.GetDB()
		if err := db.Model(&entity.AuditLog{}).Where("id = ?", id).Updates(map[string]interface{}{
			"tx_hash":        txHash,
			"author_address": authorAddress,
		}).Error; err != nil {
			fmt.Printf("Failed to update audit log with tx_hash/author: %v\n", err)
		}
	}(auditLog.ID, auditLog.BlockchainHash)

	return nil
}

// RecordActivity records user activity without response data
func (s *auditLogService) RecordActivity(userID uint64, method string, endpoint string, description string, requestBody interface{}, ipAddress string, userAgent string) error {
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
		Description: description,
		RequestBody: requestJSON,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}

	return s.createAuditLogWithBlockchain(auditLog)
}

// RecordActivityWithResponse records user activity with response data
func (s *auditLogService) RecordActivityWithResponse(userID uint64, method string, endpoint string, description string, requestBody interface{}, responseCode int, responseBody interface{}, ipAddress string, userAgent string) error {
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
		Description:  description,
		RequestBody:  requestJSON,
		ResponseCode: responseCode,
		ResponseBody: responseJSON,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}

	return s.createAuditLogWithBlockchain(auditLog)
}

// RecordActivityWithError records user activity with error information
func (s *auditLogService) RecordActivityWithError(userID uint64, method string, endpoint string, description string, requestBody interface{}, responseCode int, errorMessage string, ipAddress string, userAgent string) error {
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
		Description:  description,
		RequestBody:  requestJSON,
		ResponseCode: responseCode,
		ErrorMessage: &errorMessage,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}

	return s.createAuditLogWithBlockchain(auditLog)
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

// GetAuditLogStats retrieves statistics for audit logs
func (s *auditLogService) GetAuditLogStats() (map[string]interface{}, error) {
	totalLogs, err := s.auditLogRepository.CountAllAuditLogs()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_logs":       totalLogs,
		"contract_address": os.Getenv("STELLAR_CONTRACT_ADDRESS"),
		"network":          os.Getenv("STELLAR_NETWORK"),
	}, nil
}
