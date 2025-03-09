package entity

import "time"

type Score struct {
	ID        uint64  `gorm:"primary_key:auto_increment"`
	Score     float64 `gorm:"not null"`
	Comment   string  `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time

	SubmissionID uint64     `gorm:"not null"`
	Submission   Submission `gorm:"foreignKey:SubmissionID"`

	JudgeID uint64 `gorm:"not null"`
	Judge   Judge  `gorm:"foreignKey:JudgeID"`
}
