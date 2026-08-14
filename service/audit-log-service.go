package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gcw/entity"
	"gcw/repository"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/sahilm/fuzzy"
)

type auditLogService struct {
	auditLogRepository repository.AuditLogRepository
	userRepository     repository.UserRepository
	stellarService     StellarService
}

type AuditLogService interface {
	RecordActivity(userID uint64, userEmail string, targetEmail string, targetName string, targetResourceID uint64, method string, endpoint string, description string, requestBody interface{}, ipAddress string, userAgent string) error
	RecordActivityWithResponse(userID uint64, userEmail string, targetEmail string, targetName string, targetResourceID uint64, method string, endpoint string, description string, requestBody interface{}, responseCode int, responseBody interface{}, ipAddress string, userAgent string) error
	RecordActivityWithError(userID uint64, userEmail string, targetEmail string, targetName string, targetResourceID uint64, method string, endpoint string, description string, requestBody interface{}, responseCode int, errorMessage string, ipAddress string, userAgent string) error
	GetUserActivityLogs(userID uint64, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetEndpointActivityLogs(endpoint string, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetActivityLogsByDateRange(startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetUserActivityLogsByDateRange(userID uint64, startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetAllActivityLogs(limit int, offset int, role string, query string) ([]*entity.AuditLog, int64, error)
	GetBlockchainLogs(limit int, offset int, query string) ([]*entity.BlockchainLog, int64, error)
	GetAuditLogStats() (map[string]interface{}, error)
	GetUserIDByEmail(email string) (uint64, error)
}

func NewAuditLogService(repo repository.AuditLogRepository, userRepo repository.UserRepository, stellar StellarService) AuditLogService {
	return &auditLogService{
		auditLogRepository: repo,
		userRepository:     userRepo,
		stellarService:     stellar,
	}
}

// calculateHash returns SHA256 hash of a string
func calculateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// GetUserIDByEmail retrieves a User ID based on their email
func (s *auditLogService) GetUserIDByEmail(email string) (uint64, error) {
	var user entity.User
	err := s.userRepository.FindByEmail(email, &user)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// getStageChangeDescription extracts team name and stage from request body.
// Returns the blockchain description and true if this is a stage change, or empty string and false otherwise.
func getStageChangeDescription(auditLog *entity.AuditLog) (string, bool) {
	if auditLog.Method != "PUT" && auditLog.Method != "PATCH" {
		return "", false
	}
	if !strings.Contains(auditLog.Endpoint, "/dashboard/") {
		return "", false
	}
	if len(auditLog.RequestBody) == 0 {
		return "", false
	}

	var body map[string]interface{}
	if err := json.Unmarshal(auditLog.RequestBody, &body); err != nil {
		return "", false
	}

	stageVal, hasStage := body["stage"]
	if !hasStage {
		return "", false
	}

	stageName, _ := stageVal.(string)
	if stageName == "" {
		stageName = "unknown"
	}

	teamName := "unknown"
	if nama, ok := body["nama_tim"]; ok {
		if s, ok := nama.(string); ok && s != "" {
			teamName = s
		}
	}

	return fmt.Sprintf("Tim %s masuk ke stage %s", teamName, stageName), true
}

// createAuditLog saves the audit log to the database.
// If the log is a stage change, it also submits a SHA256 hash of the description to the Stellar blockchain.
func (s *auditLogService) createAuditLog(auditLog *entity.AuditLog) error {
	err := s.auditLogRepository.Create(auditLog)
	if err != nil {
		return err
	}

	// Submit to Stellar blockchain only for stage changes
	if blockchainDesc, ok := getStageChangeDescription(auditLog); ok {
		descHash := calculateHash(blockchainDesc)

		blockchainLog := &entity.BlockchainLog{
			AuditLogID:      auditLog.ID,
			Description:     blockchainDesc,
			DescriptionHash: descHash,
		}

		err = s.auditLogRepository.CreateBlockchainLog(blockchainLog)
		if err != nil {
			fmt.Printf("Failed to create blockchain log entry: %v\n", err)
			return err
		}

		go func(auditID uint64, blockchainLogID uint64, hash string) {
			txHash, authorAddress, err := s.stellarService.SendAuditHash(hash)
			db := s.auditLogRepository.GetDB()

			if err != nil {
				fmt.Printf("Stellar background task error: %v\n", err)
				if authorAddress != "" {
					db.Model(&entity.BlockchainLog{}).Where("id = ?", blockchainLogID).Update("author_address", authorAddress)
				}
				return
			}

			if err := db.Model(&entity.BlockchainLog{}).Where("id = ?", blockchainLogID).Updates(map[string]interface{}{
				"tx_hash":        txHash,
				"author_address": authorAddress,
			}).Error; err != nil {
				fmt.Printf("Failed to update blockchain log with tx_hash/author: %v\n", err)
			}
		}(auditLog.ID, blockchainLog.ID, descHash)
	}

	return nil
}

// RecordActivity records user activity without response data
func (s *auditLogService) RecordActivity(userID uint64, userEmail string, targetEmail string, targetName string, targetResourceID uint64, method string, endpoint string, description string, requestBody interface{}, ipAddress string, userAgent string) error {
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

	return s.createAuditLog(auditLog)
}

// RecordActivityWithResponse records user activity with response data
func (s *auditLogService) RecordActivityWithResponse(userID uint64, userEmail string, targetEmail string, targetName string, targetResourceID uint64, method string, endpoint string, description string, requestBody interface{}, responseCode int, responseBody interface{}, ipAddress string, userAgent string) error {
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

	return s.createAuditLog(auditLog)
}

// RecordActivityWithError records user activity with error information
func (s *auditLogService) RecordActivityWithError(userID uint64, userEmail string, targetEmail string, targetName string, targetResourceID uint64, method string, endpoint string, description string, requestBody interface{}, responseCode int, errorMessage string, ipAddress string, userAgent string) error {
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

	return s.createAuditLog(auditLog)
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

// GetBlockchainLogs retrieves all blockchain log entries (both confirmed and pending)
func (s *auditLogService) GetBlockchainLogs(limit int, offset int, query string) ([]*entity.BlockchainLog, int64, error) {
	db := s.auditLogRepository.GetDB()
	var logs []*entity.BlockchainLog
	var total int64

	// Build a sub-query to get matching IDs (avoids Preload + Joins conflict)
	idsQuery := db.Model(&entity.BlockchainLog{}).Select("blockchain_logs.id")

	if query != "" {
		idsQuery = idsQuery.
			Joins("JOIN audit_logs ON audit_logs.id = blockchain_logs.audit_log_id").
			Where("blockchain_logs.description ILIKE ? OR audit_logs.description ILIKE ? OR blockchain_logs.tx_hash ILIKE ? OR blockchain_logs.description_hash ILIKE ?",
				"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	var matchingIDs []uint64
	if err := idsQuery.Pluck("blockchain_logs.id", &matchingIDs).Error; err != nil {
		return nil, 0, err
	}
	total = int64(len(matchingIDs))

	if total == 0 {
		return []*entity.BlockchainLog{}, 0, nil
	}

	err := db.Model(&entity.BlockchainLog{}).
		Where("id IN ?", matchingIDs).
		Preload("AuditLog").Preload("AuditLog.User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAllActivityLogs retrieves all activity logs with pagination, optional role filter, and smart fuzzy search
func (s *auditLogService) GetAllActivityLogs(limit int, offset int, role string, query string) ([]*entity.AuditLog, int64, error) {
	if query != "" {
		// Fetch ALL logs filtered by role for fuzzy matching
		// We use a high limit to get all relevant records for in-memory fuzzy search
		allLogs, _, err := s.auditLogRepository.FindAll(100000, 0, role, "")
		if err != nil {
			return nil, 0, err
		}

		// Perform fuzzy search
		matches := fuzzy.FindFrom(query, auditLogSource(allLogs))

		// Sort by relevance score
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Score > matches[j].Score
		})

		var filteredData []*entity.AuditLog
		for _, match := range matches {
			filteredData = append(filteredData, allLogs[match.Index])
		}

		total := int64(len(filteredData))

		// Manual pagination
		start := offset
		if start > len(filteredData) {
			return []*entity.AuditLog{}, total, nil
		}

		end := start + limit
		if end > len(filteredData) {
			end = len(filteredData)
		}

		return filteredData[start:end], total, nil
	}

	return s.auditLogRepository.FindAll(limit, offset, role, "")
}

// Fuzzy Search Sources for Audit Logs
type auditLogSource []*entity.AuditLog

func (s auditLogSource) String(i int) string {
	// Search against description and user email
	return s[i].Description + " " + s[i].User.Email
}

func (s auditLogSource) Len() int {
	return len(s)
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
