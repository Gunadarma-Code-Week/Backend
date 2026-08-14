package entity

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID           uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64          `gorm:"not null; index" json:"user_id"`
	User         User            `gorm:"foreignKey:UserID" json:"user"`
	Method       string          `gorm:"varchar(10); not null" json:"method"` // POST, PUT, DELETE, PATCH (GET not logged)
	Endpoint     string          `gorm:"varchar(255); not null; index" json:"endpoint"`
	Description  string          `gorm:"type:text" json:"description"` // Human-readable description of the action
	RequestBody  json.RawMessage `gorm:"type:jsonb" json:"request_body"`
	ResponseCode int             `gorm:"default:0" json:"response_code"`
	ResponseBody json.RawMessage `gorm:"type:jsonb" json:"response_body"`
	IPAddress    string          `gorm:"varchar(45)" json:"ip_address"`
	UserAgent    string          `gorm:"type:text" json:"user_agent"`
	ErrorMessage *string         `gorm:"type:text" json:"error_message"`
	CreatedAt    time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
