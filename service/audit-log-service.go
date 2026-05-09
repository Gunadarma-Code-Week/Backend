package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
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
	stellarService     StellarService
	secretKey          []byte
}

type AuditLogService interface {
	RecordActivity(userID uint64, userEmail string, method string, endpoint string, description string, requestBody interface{}, ipAddress string, userAgent string) error
	RecordActivityWithResponse(userID uint64, userEmail string, method string, endpoint string, description string, requestBody interface{}, responseCode int, responseBody interface{}, ipAddress string, userAgent string) error
	RecordActivityWithError(userID uint64, userEmail string, method string, endpoint string, description string, requestBody interface{}, responseCode int, errorMessage string, ipAddress string, userAgent string) error
	GetUserActivityLogs(userID uint64, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetEndpointActivityLogs(endpoint string, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetActivityLogsByDateRange(startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetUserActivityLogsByDateRange(userID uint64, startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error)
	GetAllActivityLogs(limit int, offset int, role string, query string) ([]*entity.AuditLog, int64, error)
	GetAuditLogStats() (map[string]interface{}, error)
}

func NewAuditLogService(repo repository.AuditLogRepository, stellar StellarService) AuditLogService {
	// Load audit_private.pem as secret key for HMAC
	auditKeyPath := os.Getenv("AUDIT_PRIVATE_KEY_PATH")
	if auditKeyPath == "" {
		auditKeyPath = "keys/audit_private.pem"
	}

	secretKey, err := os.ReadFile(auditKeyPath)
	if err != nil {
		fmt.Printf("Warning: Failed to load %s for audit log HMAC: %v. Using empty key.\n", auditKeyPath, err)
		secretKey = []byte("")
	}

	return &auditLogService{
		auditLogRepository: repo,
		stellarService:     stellar,
		secretKey:          secretKey,
	}
}

// calculateHash returns SHA256 hash of a string
func (s *auditLogService) calculateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// calculateHMAC returns HMAC-SHA256 hash of a string using the secret key
func (s *auditLogService) calculateHMAC(data string) string {
	h := hmac.New(sha256.New, s.secretKey)
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// calculateAESGCMHash encrypts data with AES-256-GCM then hashes the result with SHA-256
func (s *auditLogService) calculateAESGCMHash(data string) string {
	// 1. Derive 32-byte key for AES-256 from secretKey
	key := sha256.Sum256(s.secretKey)
	
	block, err := aes.NewCipher(key[:])
	if err != nil {
		// Fallback to HMAC if AES fails (should not happen with valid key)
		return s.calculateHMAC(data)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return s.calculateHMAC(data)
	}

	// 2. Create a deterministic nonce from the data itself (so hash is reproducible)
	// We take first 12 bytes of SHA256(data) as nonce
	nonceHash := sha256.Sum256([]byte(data))
	nonce := nonceHash[:gcm.NonceSize()]

	// 3. Encrypt data
	ciphertext := gcm.Seal(nil, nonce, []byte(data), nil)
	
	// 4. Hash the combined result (nonce + ciphertext) with SHA-256 to get 64 chars
	finalHash := sha256.Sum256(append(nonce, ciphertext...))
	
	return fmt.Sprintf("%x", finalHash)
}

// createAuditLogWithBlockchain adds blockchain hashes and saves the audit log
func (s *auditLogService) createAuditLogWithBlockchain(auditLog *entity.AuditLog, userEmail string) error {
	// Ensure CreatedAt is set before hashing for consistency
	if auditLog.CreatedAt.IsZero() {
		auditLog.CreatedAt = time.Now()
	}

	// Create a description for hashing that replaces email with User ID for privacy
	blockchainDescription := auditLog.Description
	if userEmail != "" {
		blockchainDescription = strings.ReplaceAll(blockchainDescription, userEmail, fmt.Sprintf("User (ID: %d)", auditLog.UserID))
	}

	// Create a dictionary for hashing that includes metadata
	metadata := map[string]string{
		"timestamp":      auditLog.CreatedAt.Format(time.RFC3339),
		"deskripsi":      blockchainDescription,
		"ip":             auditLog.IPAddress,
		"method":         auditLog.Method,
		"endpoint":       auditLog.Endpoint,
		"status_code":    fmt.Sprintf("%d", auditLog.ResponseCode),
		"wallet_address": s.stellarService.GetPublicKey(),
	}
	metadataJSON, _ := json.Marshal(metadata)

	// 1. Calculate AES-GCM-SHA256 hash of metadata dictionary
	descriptionHash := s.calculateAESGCMHash(string(metadataJSON))
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

	// 3. Calculate current blockchain hash using AES-GCM-SHA256: 
	// Process (metadataJSON + lastBlockchainHash) through AES-GCM then hash with SHA256
	currentBlockchainHash := s.calculateAESGCMHash(string(metadataJSON) + lastBlockchainHash)
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
func (s *auditLogService) RecordActivity(userID uint64, userEmail string, method string, endpoint string, description string, requestBody interface{}, ipAddress string, userAgent string) error {
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

	return s.createAuditLogWithBlockchain(auditLog, userEmail)
}

// RecordActivityWithResponse records user activity with response data
func (s *auditLogService) RecordActivityWithResponse(userID uint64, userEmail string, method string, endpoint string, description string, requestBody interface{}, responseCode int, responseBody interface{}, ipAddress string, userAgent string) error {
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

	return s.createAuditLogWithBlockchain(auditLog, userEmail)
}

// RecordActivityWithError records user activity with error information
func (s *auditLogService) RecordActivityWithError(userID uint64, userEmail string, method string, endpoint string, description string, requestBody interface{}, responseCode int, errorMessage string, ipAddress string, userAgent string) error {
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

	return s.createAuditLogWithBlockchain(auditLog, userEmail)
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
