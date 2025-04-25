package entity

import "time"

type Team struct {
	ID_Team        uint64 `gorm:"primary_key:auto_increment"`
	TeamName       string `gorm:"varchar(255); not null"`
	Supervisor     string `gorm:"varchar(255); not null"`
	SupervisorNIDN string `gorm:"varchar(255); not null"`
	JoinCode       string `gorm:"varchar(255); not null"`
	KomitmenFee    string `gorm:"varchar(255)"`

	Event string `gorm:"varchar(255); not null"`

	ID_LeadTeam uint64 `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
