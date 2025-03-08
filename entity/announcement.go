package entity

import "time"

type Announcement struct {
	ID        uint64 `gorm:"primary_key:auto_increment"`
	EventID   uint64 `gorm:"not null"`
	Title     string `gorm:"type:varchar(255);not null"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Event Event `gorm:"foreignKey:EventID"`
}
