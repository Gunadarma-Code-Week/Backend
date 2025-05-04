package entity

import (
	"time"
)

type User struct {
	ID    uint64 `gorm:"primary_key:auto_increment"`
	Email string `gorm:"varchar(255); not null"`

	Name            string     `gorm:"varchar(55); not null"`
	Gender          *string    `gorm:"varchar(55);"`
	NIM             *string    `gorm:"varchar(55);"`
	BirthPlace      *string    `gorm:"varchar(55);"`
	BirthDate       *time.Time `gorm:"date"`
	Institusi       string     `gorm:"varchar(55);"`
	Phone           string     `gorm:"type:varchar(16)"`
	DokumenFilename string     `gorm:"type:varchar(255)"`
	SocMedDocument  string     `gorm:"type:varchar(255)"`
	Jenjang         string     `gorm:"type:varchar(120)"`
	Major           string     `gorm:"type:varchar(120)"`
	ProfilePicture  string     `gorm:"type:varchar(255)"`

	ProfileHasUpdated bool `gorm:"bool; default:false"`
	DataHasVerified   bool `gorm:"bool; default:false"`

	// RoleID uint64   `gorm:"not null"` // Kolom untuk join ke tabel User_Role
	Role string `gorm:"varchar(255); not null; default:'user'"` // user, admin, superadmin

	IDTeam *uint64
	Team   Team `gorm:"foreignKey:IDTeam;references:ID_Team"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
