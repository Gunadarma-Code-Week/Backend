package entity

import "time"

type Team struct {
	ID_Team        uint64 `gorm:"primaryKey;autoIncrement"`
	TeamName       string `gorm:"type:varchar(255); not null"`
	Supervisor     string `gorm:"type:varchar(255)"`
	SupervisorNIDN string `gorm:"type:varchar(255)"`
	JoinCode       string `gorm:"type:varchar(255); not null"`
	KomitmenFee    string `gorm:"type:varchar(255)" json:"komitmen_fee"`
	Event          string `gorm:"type:varchar(255); not null" json:"event"`

	PaymentStatus string `gorm:"type:varchar(50); default:'pending'" json:"payment_status"`
	OrderID       string `gorm:"type:varchar(255);uniqueIndex" json:"order_id"`
	QRString      string `gorm:"type:text" json:"qr_string"`

	ID_LeadTeam uint64 `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
