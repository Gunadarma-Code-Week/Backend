package entity

import "time"

type Team struct {
	ID_Team        uint64 `gorm:"primary_key:auto_increment"`
	TeamName       string `gorm:"varchar(255); not null"`
	Supervisor     string `gorm:"varchar(255)"`
	SupervisorNIDN string `gorm:"varchar(255)"`
	JoinCode       string `gorm:"varchar(255); not null"`
	KomitmenFee    string `gorm:"varchar(255)"`

	Event string `gorm:"varchar(255); not null"`

	PaymentStatus string `gorm:"varchar(50); default:'pending'"`
	OrderID       string `gorm:"varchar(255); unique"`
	QRString      string `gorm:"text"`

	ID_LeadTeam uint64 `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
