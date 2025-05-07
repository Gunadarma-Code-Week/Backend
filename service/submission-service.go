package service

import (
	"fmt"
	"gcw/dto"
	"gcw/entity"

	"gorm.io/gorm"
)

// Interface
type SubmissionService interface {
	Create(string, string, dto.RequestHackathon) (entity.HackathonTeam, error)
	Get(string) (dto.HackatonStageStatus, error)
}

// Concrete implementation
type submissionService struct {
	db *gorm.DB
}

// Constructor
func NewSubmissionService(db *gorm.DB) SubmissionService {
	return &submissionService{db: db}
}

func (s *submissionService) Create(join_code, stage string, submissionDTO dto.RequestHackathon) (entity.HackathonTeam, error) {
	var submission entity.HackathonTeam
	var team entity.Team

	if err := s.db.Where("join_code = ?", join_code).First(&team).Error; err != nil {
		return entity.HackathonTeam{}, err
	}

	if err := s.db.Where("id_team = ?", team.ID_Team).First(&submission).Error; err != nil {
		return entity.HackathonTeam{}, fmt.Errorf("team tidak terdaftar")
	}

	switch stage {
	case "stage1":
		submission.ProposalUrl = submissionDTO.LinkDrive
	case "stage2":
		submission.PitchDeckUrl = submissionDTO.LinkDrive
	case "final":
		submission.GithubProjectUrl = submissionDTO.LinkDrive
	}

	if err := s.db.Save(&submission).Error; err != nil {
		return entity.HackathonTeam{}, err
	}

	return submission, nil
}

func (s *submissionService) Get(join_code string) (dto.HackatonStageStatus, error) {
	var submission entity.HackathonTeam
	var team entity.Team

	if err := s.db.Where("join_code = ?", join_code).First(&team).Error; err != nil {
		return dto.HackatonStageStatus{}, err
	}

	if err := s.db.Where("id_team = ?", team.ID_Team).First(&submission).Error; err != nil {
		return dto.HackatonStageStatus{}, err
	}

	var status dto.HackatonStageStatus

	if submission.ProposalUrl != "" {
		status.Stage1 = true
		status.Stage1Url = submission.ProposalUrl
	}
	if submission.PitchDeckUrl != "" {
		status.Stage2 = true
		status.Stage2Url = submission.PitchDeckUrl
	}
	if submission.GithubProjectUrl != "" {
		status.Final = true
		status.FinalUrl = submission.GithubProjectUrl
	}

	return status, nil
}
