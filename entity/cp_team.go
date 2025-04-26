package entity

import "time"

type CPTeam struct {
	ID_CPTeam        uint64 `gorm:"primary_key:auto_increment"`
	Stage            string `gorm:"varchar(255); not null"`
	Status           string `gorm:"varchar(255); not null"`
	DomjudgeUsername string `gorm:"varchar(255); not null"`
	DomjudgePassword string `gorm:"varchar(255); not null"`

	IDTeam uint64 `gorm:"not null"`
	Team   Team   `gorm:"foreignKey:IDTeam"`

	IsDeleted     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
