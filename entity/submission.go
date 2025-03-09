package entity

import "time"

type Submission struct {
	ID_Submission uint64 `gorm:"primary_key:auto_increment"`
	FileURL       string `gorm:"type:text"`
	TotalScore    *float64

	EventID uint64 `gorm:"not null"`
	Event   Event  `gorm:"foreignKey:EventID"`

	TeamID uint64 `gorm:"not null"`
	Team   Team   `gorm:"foreignKey:TeamID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
