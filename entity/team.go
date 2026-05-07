package entity

import "time"

type Team struct {
	ID_Team        uint64 `gorm:"primaryKey;autoIncrement"`
	TeamName       string `gorm:"type:varchar(255); not null"`
	Supervisor     string `gorm:"type:varchar(255)"`
	SupervisorNIDN string `gorm:"type:varchar(255)"`
	JoinCode       string `gorm:"type:varchar(255); not null"`
	KomitmenFee    string `gorm:"type:varchar(255)"`

	Event string `gorm:"type:varchar(255); not null"`

	PaymentStatus string `gorm:"type:varchar(50); default:'pending'"`
	OrderID       string `gorm:"type:varchar(255);uniqueIndex"`
	QRString      string `gorm:"type:text"`

	ID_LeadTeam uint64 `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
