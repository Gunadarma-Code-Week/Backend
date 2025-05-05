package dto

type CpDetailDto struct {
	Stage            string `gorm:"varchar(255); not null"`
	Status           string `gorm:"varchar(255); not null"`
	DomjudgeUsername string `gorm:"varchar(255); not null"`
	DomjudgePassword string `gorm:"varchar(255); not null"`
}
