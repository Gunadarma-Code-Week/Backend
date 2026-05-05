package entity

import "time"

type CTFTeam struct {
	ID_CTFTeam uint64 `gorm:"primary_key:auto_increment"`
	Stage      string `gorm:"varchar(255); not null"`
	Status           string `gorm:"varchar(255); not null"`
	DomjudgeUsername string `gorm:"varchar(255)"`
	DomjudgePassword string `gorm:"varchar(255)"`

	IDTeam uint64 `gorm:"not null"`
	Team   Team   `gorm:"foreignKey:IDTeam"`

	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
