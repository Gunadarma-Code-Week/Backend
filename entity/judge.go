package entity

import "time"

type Judge struct {
	ID        uint64 `gorm:"primary_key:auto_increment"`
	Name      string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"type:varchar(255);uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
