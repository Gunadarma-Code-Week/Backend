package entity

import "time"

type Score struct {
	ID           uint64  `gorm:"primary_key:auto_increment"`
	SubmissionID uint64  `gorm:"not null"`
	JudgeID      uint64  `gorm:"not null"`
	Score        float64 `gorm:"not null"`
	Comment      string  `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Submission Submission `gorm:"foreignKey:SubmissionID"`
	Judge      Judge      `gorm:"foreignKey:JudgeID"`
}
