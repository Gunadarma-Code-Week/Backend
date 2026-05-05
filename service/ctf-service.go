package service

import (
	"gcw/dto"
	"gcw/entity"

	"gorm.io/gorm"
)

type CtfService struct {
	db *gorm.DB
}

func NewCtfService(db *gorm.DB) *CtfService {
	return &CtfService{db: db}
}

func (s *CtfService) Get(join_code string) (dto.CTFDetailDto, error) {
	var team entity.Team
	var ctfTeam entity.CTFTeam

	if err := s.db.Where("join_code = ?", join_code).First(&team).Error; err != nil {
		return dto.CTFDetailDto{}, err
	}

	if err := s.db.Where("id_team = ?", team.ID_Team).First(&ctfTeam).Error; err != nil {
		return dto.CTFDetailDto{}, err
	}

	return dto.CTFDetailDto{
		Stage:            ctfTeam.Stage,
		Status:           ctfTeam.Status,
		DomjudgeUsername: ctfTeam.DomjudgeUsername,
		DomjudgePassword: ctfTeam.DomjudgePassword,
	}, nil
}
