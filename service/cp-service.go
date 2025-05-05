package service

import (
	"gcw/dto"
	"gcw/entity"

	"gorm.io/gorm"
)

// Concrete implementation
type CpService struct {
	db *gorm.DB
}

// Constructor
func NewCpService(db *gorm.DB) *CpService {
	return &CpService{db: db}
}

func (s *CpService) Get(join_code string) (dto.CpDetailDto, error) {
	var team entity.Team
	var cpTeam entity.CPTeam

	if err := s.db.Where("join_code = ?", join_code).First(&team).Error; err != nil {
		return dto.CpDetailDto{}, err
	}

	if err := s.db.Where("id_team = ?", team.ID_Team).First(&cpTeam).Error; err != nil {
		return dto.CpDetailDto{}, err
	}

	var status dto.CpDetailDto

	status.Stage = cpTeam.Stage
	status.Status = cpTeam.Status
	status.DomjudgeUsername = cpTeam.DomjudgeUsername
	status.DomjudgePassword = cpTeam.DomjudgePassword

	return status, nil
}
