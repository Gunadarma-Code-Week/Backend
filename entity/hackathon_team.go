package entity

import "time"

type HackathonTeam struct {
	ID_HackathonTeam uint64 `gorm:"primaryKey;autoIncrement"`
	Stage            string `gorm:"varchar(255); not null"`
	Status           string `gorm:"varchar(255); not null"`
	ProposalUrl      string `gorm:"varchar(255); null"`
	GithubProjectUrl string `gorm:"varchar(255); null"`
	PitchDeckUrl     string `gorm:"varchar(255); null"`

	IDTeam uint64
	Team   Team `gorm:"foreignKey:IDTeam"`

	IsDeleted bool
	DeletionReason string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
