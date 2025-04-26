package entity

import "time"

type HackathonTeam struct {
	ID_HackathonTeam uint64 `gorm:"primary_key:auto_increment"`
	Stage            string `gorm:"varchar(255); not null"`
	Status           string `gorm:"varchar(255); not null"`
	ProposalUrl      string `gorm:"varchar(255); null"`
	GithubProjectUrl string `gorm:"varchar(255); null"`
	PitchDeckUrl     string `gorm:"varchar(255); null"`

	IDTeam uint64
	Team   Team `gorm:"foreignKey:IDTeam"`

	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
