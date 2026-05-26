package entity

type Timeline struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Category    string `gorm:"type:varchar(50);not null" json:"category"` // "home", "hackathon", "cp", "ctf"
	OrderIndex  int    `gorm:"not null;default:0" json:"order_index"`
	Date        string `gorm:"type:varchar(100);not null" json:"date"`
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Events      string `gorm:"type:text" json:"events"` // JSON array string
}
