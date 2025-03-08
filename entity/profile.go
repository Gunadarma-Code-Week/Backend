package entity

import (
	"time"
)

type Profile struct {
	ID_Profile uint64    `gorm:"primary_key:auto_increment"`
	Name       string    `gorm:"varchar(55); not null"`
	Gender     string    `gorm:"varchar(55); not null"`
	NIM        uint64    `gorm:"not null"`
	Age        uint64    `gorm:"not null"`
	BirthPlace string    `gorm:"varchar(55); not null"`
	BirthDate  time.Time `gorm:"not null"`
	Institusi  string    `gorm:"varchar(55); not null"`

	UserID uint64   `gorm:"not null"` // Kolom untuk join ke tabel User_Role
	User   UserRole `gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
