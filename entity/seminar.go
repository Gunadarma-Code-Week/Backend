package entity

import "time"

type Seminar struct {
	ID_Seminar uint64 `gorm:"primary_key:auto_increment"`
	Name       string `gorm:"varchar(255); not null"`

	ID_Event uint64 `gorm:"not null"`
	Event    Event  `gorm:"foreignKey:ID_Event;references:ID_Event"`

	Users []User `gorm:"many2many:users_member"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
