package entity

import "time"

type NewsLetter struct {
	ID_NewsLetter uint64 `gorm:"primary_key;auto_increment" json:"id_news_letter"`
	Title         string `gorm:"type:varchar(255);not null" json:"title"`
	NewsLetter    string `gorm:"type:text;not null" json:"news_letter"`
	BaseImage     string `gorm:"type:varchar(255)" json:"base_image"`
	ID_Admin      uint64 `gorm:"not null" json:"id_admin"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
