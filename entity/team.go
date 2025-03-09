package entity

import "time"

type Team struct {
	ID_Team        uint64 `gorm:"primary_key:auto_increment"`
	Name           string `gorm:"varchar(255); not null"`
	Supervisor     string `gorm:"varchar(255); not null"`
	SupervisorNIDN string `gorm:"varchar(255); not null"`

	ID_LeadTeam uint64 `gorm:"not null"`
	UserLead    User   `gorm:"foreignKey:ID_LeadTeam"`

	ID_Event uint64 `gorm:"not null"`
	Event    Event  `gorm:"foreignKey:ID_Event"`

	Users []User `gorm:"many2many:team_users"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
