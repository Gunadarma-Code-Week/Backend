package entity

import "time"

type UserRole struct {
	ID_Role   uint64 `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"varchar(255); not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
