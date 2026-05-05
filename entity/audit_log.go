package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID           uint64          `gorm:"primary_key:auto_increment"`
	UserID       uint64          `gorm:"not null; index"`
	Method       string          `gorm:"varchar(10); not null"` // POST, PUT, DELETE, PATCH (GET not logged)
	Endpoint     string          `gorm:"varchar(255); not null; index"`
	RequestBody  json.RawMessage `gorm:"type:jsonb"`
	ResponseCode int             `gorm:"default:0"`
	ResponseBody json.RawMessage `gorm:"type:jsonb"`
	IPAddress    string          `gorm:"varchar(45)"`
	UserAgent    string          `gorm:"type:text"`
	ErrorMessage *string         `gorm:"type:text"`
	CreatedAt    time.Time       `gorm:"autoCreateTime"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime"`
}

// Scan implements sql.Scanner interface
func (a *AuditLog) Scan(value interface{}) error {
	return nil
}

// Value implements driver.Valuer interface
func (a AuditLog) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
