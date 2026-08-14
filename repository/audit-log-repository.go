package repository

import (
	"gcw/entity"
	"time"

	"gorm.io/gorm"
)

type auditLogRepository struct {
	DB *gorm.DB
}

type AuditLogRepository interface {
	Create(auditLog *entity.AuditLog) error
	CreateBlockchainLog(blockchainLog *entity.BlockchainLog) error
	FindByUserID(userID uint64, limit int, offset int) ([]*entity.AuditLog, int64, error)
	FindByEndpoint(endpoint string, limit int, offset int) ([]*entity.AuditLog, int64, error)
	FindByDateRange(startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error)
	FindByUserIDAndDateRange(userID uint64, startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error)
	FindAll(limit int, offset int, role string, query string) ([]*entity.AuditLog, int64, error)
	CountAllAuditLogs() (int64, error)
	GetLastAuditLog() (*entity.AuditLog, error)
	GetDB() *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{
		DB: db,
	}
}

func (r *auditLogRepository) GetDB() *gorm.DB {
	return r.DB
}

// Create a new audit log record
func (r *auditLogRepository) Create(auditLog *entity.AuditLog) error {
	res := r.DB.Create(auditLog)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// Create a new blockchain log record
func (r *auditLogRepository) CreateBlockchainLog(blockchainLog *entity.BlockchainLog) error {
	res := r.DB.Create(blockchainLog)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// FindByUserID retrieves audit logs for a specific user
func (r *auditLogRepository) FindByUserID(userID uint64, limit int, offset int) ([]*entity.AuditLog, int64, error) {
	var auditLogs []*entity.AuditLog
	var total int64

	res := r.DB.Preload("User").Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs)

	if res.Error != nil {
		return nil, 0, res.Error
	}

	r.DB.Model(&entity.AuditLog{}).Where("user_id = ?", userID).Count(&total)

	return auditLogs, total, nil
}

// FindByEndpoint retrieves audit logs for a specific endpoint
func (r *auditLogRepository) FindByEndpoint(endpoint string, limit int, offset int) ([]*entity.AuditLog, int64, error) {
	var auditLogs []*entity.AuditLog
	var total int64

	res := r.DB.Preload("User").Where("endpoint = ?", endpoint).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs)

	if res.Error != nil {
		return nil, 0, res.Error
	}

	r.DB.Model(&entity.AuditLog{}).Where("endpoint = ?", endpoint).Count(&total)

	return auditLogs, total, nil
}

// FindByDateRange retrieves audit logs within a date range
func (r *auditLogRepository) FindByDateRange(startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error) {
	var auditLogs []*entity.AuditLog
	var total int64

	res := r.DB.Preload("User").Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs)

	if res.Error != nil {
		return nil, 0, res.Error
	}

	r.DB.Model(&entity.AuditLog{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&total)

	return auditLogs, total, nil
}

// FindByUserIDAndDateRange retrieves audit logs for a user within a date range
func (r *auditLogRepository) FindByUserIDAndDateRange(userID uint64, startDate time.Time, endDate time.Time, limit int, offset int) ([]*entity.AuditLog, int64, error) {
	var auditLogs []*entity.AuditLog
	var total int64

	res := r.DB.Preload("User").Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs)

	if res.Error != nil {
		return nil, 0, res.Error
	}

	r.DB.Model(&entity.AuditLog{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).Count(&total)

	return auditLogs, total, nil
}

// FindAll retrieves all audit logs with pagination, role filter, and search query
func (r *auditLogRepository) FindAll(limit int, offset int, role string, searchQuery string) ([]*entity.AuditLog, int64, error) {
	var auditLogs []*entity.AuditLog
	var total int64

	dbQuery := r.DB.Model(&entity.AuditLog{}).Preload("User").
		Joins("LEFT JOIN users ON users.id = audit_logs.user_id")

	if role != "" {
		dbQuery = dbQuery.Where("users.role = ?", role)
	}

	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		dbQuery = dbQuery.Where("(audit_logs.description ILIKE ? OR users.email ILIKE ?)", searchPattern, searchPattern)
	}

	res := dbQuery.Order("audit_logs.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs)

	if res.Error != nil {
		return nil, 0, res.Error
	}

	countQuery := r.DB.Model(&entity.AuditLog{}).
		Joins("LEFT JOIN users ON users.id = audit_logs.user_id")

	if role != "" {
		countQuery = countQuery.Where("users.role = ?", role)
	}

	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		countQuery = countQuery.Where("(audit_logs.description ILIKE ? OR users.email ILIKE ?)", searchPattern, searchPattern)
	}

	countQuery.Count(&total)

	return auditLogs, total, nil
}

func (r *auditLogRepository) CountAllAuditLogs() (int64, error) {
	var total int64
	if err := r.DB.Model(&entity.AuditLog{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// GetLastAuditLog retrieves the most recent audit log
func (r *auditLogRepository) GetLastAuditLog() (*entity.AuditLog, error) {
	var auditLog entity.AuditLog
	res := r.DB.Order("id DESC").First(&auditLog)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, res.Error
	}
	return &auditLog, nil
}
