package entity

import "time"

type User struct {
	ID        uint64 `gorm:"primary_key:auto_increment"`
	Name      string `gorm:"varchar(255); not null"`
	Username  string `gorm:"type:varchar(255); index:username_index,type:btree; not null"`
	Password  string `gorm:"text; not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
