package entity

import "time"

type User struct {
	ID    uint64 `gorm:"primary_key:auto_increment"`
	Email string `gorm:"varchar(255); not null"`

	// RoleID uint64   `gorm:"not null"` // Kolom untuk join ke tabel User_Role
	Role *string `gorm:"varchar(255); not null; default:'user'"` // user, admin, superadmin

	IDTeam *uint64
	Team   Team `gorm:"foreignKey:IDTeam"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
