package entity

import (
	"time"
)

type BlockchainLog struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	AuditLogID      uint64    `gorm:"not null; index; unique" json:"audit_log_id"`
	AuditLog        AuditLog  `gorm:"foreignKey:AuditLogID" json:"audit_log"`
	Description     string    `gorm:"type:text" json:"description"`             // Human-readable description (e.g. "Tim X masuk ke stage Y")
	DescriptionHash string    `gorm:"type:varchar(64)" json:"description_hash"` // SHA256 of the description
	TxHash          string    `gorm:"type:varchar(66)" json:"tx_hash"`          // Stellar transaction hash
	AuthorAddress   string    `gorm:"type:varchar(56)" json:"author_address"`   // Stellar public key of the signer
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (BlockchainLog) TableName() string {
	return "blockchain_logs"
}
