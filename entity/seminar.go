package entity

import "time"

type Seminar struct {
	ID_Seminar uint64 `gorm:"primary_key:auto_increment"`
	ID_Tiket   string `gorm:"varchar(255); not null"`

	IDUser uint64 `gorm:"not null"`
	User   User   `gorm:"foreignKey:IDUser;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
