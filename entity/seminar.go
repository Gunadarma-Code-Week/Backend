package entity

import "time"

type Seminar struct {
	ID_Seminar uint64 `gorm:"primaryKey;autoIncrement"`
	ID_Tiket   string `gorm:"varchar(255); not null"`

	IDUser        uint64 `gorm:"not null"`
	User          User   `gorm:"foreignKey:IDUser;references:ID"`
	PaymentStatus string `gorm:"type:varchar(60)"`
	
	IsDeleted     bool
	DeletionReason string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
