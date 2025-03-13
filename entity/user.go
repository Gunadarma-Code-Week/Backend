package entity

import "time"

type User struct {
	ID          uint64 `gorm:"primary_key:auto_increment"`
	GoogleSubId uint64 `gorm:"not null"`
	Email       string `gorm:"varchar(255); not null"`
	Username    string `gorm:"type:varchar(255); index:username_index,type:btree; not null"`
	Password    string `gorm:"text; not null"`

	RoleID uint64   `gorm:"not null"` // Kolom untuk join ke tabel User_Role
	Role   UserRole `gorm:"foreignKey:RoleID"`

	IDTeam *uint64
	Team   Team `gorm:"foreignKey:IDTeam"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
