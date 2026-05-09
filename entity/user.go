package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Email    string `gorm:"varchar(255); not null" json:"email"`
	Password string `gorm:"type:varchar(255)" json:"-"`

	Name            string     `gorm:"varchar(55); not null" json:"name"`
	NIM             *string    `gorm:"varchar(55);" json:"nim"`
	Institusi       string     `gorm:"varchar(55);" json:"institusi"`
	Phone           string     `gorm:"type:varchar(16)" json:"phone"`
	DokumenFilename string     `gorm:"type:varchar(255)" json:"dokumen_filename"`
	SocMedDocument  string     `gorm:"type:varchar(255)" json:"soc_med_document"`
	Jenjang         string     `gorm:"type:varchar(120)" json:"jenjang"`
	Major           string     `gorm:"type:varchar(120)" json:"major"`
	ProfilePicture  string     `gorm:"type:varchar(255)" json:"profile_picture"`

	ProfileHasUpdated bool `gorm:"bool; default:false" json:"profile_has_updated"`
	DataHasVerified   bool `gorm:"bool; default:false" json:"data_has_verified"`

	Role string `gorm:"varchar(255); not null; default:'user'" json:"role"` 

	IDTeam *uint64 `json:"id_team"`
	Team   Team    `gorm:"foreignKey:IDTeam;references:ID_Team" json:"team"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	DeletionReason string    `gorm:"type:text" json:"deletion_reason"`
}
