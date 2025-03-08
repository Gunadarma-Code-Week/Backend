package entity

import "time"

type Submission struct {
	ID_Submission uint64 `gorm:"primary_key:auto_increment"`
	EventID       uint64 `gorm:"not null"`
	TeamID        uint64 `gorm:"not null"`
	FileURL       string `gorm:"type:text"`
	TotalScore    *float64

	Event Event `gorm:"foreignKey:EventID"`
	Team  Team  `gorm:"foreignKey:TeamID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
