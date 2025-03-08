package entity

import "time"

type Event struct {
	ID_Event uint64 `gorm:"primary_key:auto_increment"`
	Name     string `gorm:"varchar(255); not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
