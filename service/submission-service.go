package service

import (
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"time"

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

func parseDeadline(deadlineStr string) (time.Time, error) {
	if deadlineStr == "" {
		return time.Time{}, fmt.Errorf("empty deadline")
	}
	wib := time.FixedZone("WIB", 7*60*60)

	// 1. Try ISO8601 / RFC3339 (e.g. 2026-05-19T12:00:00+07:00)
	if t, err := time.Parse(time.RFC3339, deadlineStr); err == nil {
		return t.In(wib), nil
	}

	// 2. Try datetime-local format with seconds (e.g. 2026-05-19T12:00:00)
	if t, err := time.ParseInLocation("2006-01-02T15:04:05", deadlineStr, wib); err == nil {
		return t, nil
	}

	// 3. Try datetime-local format without seconds (e.g. 2026-05-19T12:00)
	if t, err := time.ParseInLocation("2006-01-02T15:04", deadlineStr, wib); err == nil {
		return t, nil
	}

	// 4. Try date-only format (e.g. 2026-05-19)
	if t, err := time.ParseInLocation("2006-01-02", deadlineStr, wib); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unknown format: %s", deadlineStr)
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

	// Enforce dynamic system setting block & deadline date check
	var setting entity.SystemSetting
	if err := s.db.First(&setting, 1).Error; err == nil {
		// Enforce Proposal Block & Deadline (Stage 1)
		if stage == "stage1" {
			if setting.HackathonProposalDisabled {
				return entity.HackathonTeam{}, fmt.Errorf("pengumpulan proposal (Stage 1) telah ditutup oleh administrator")
			}
			deadlineStr := setting.HackathonProposalDeadline
			if deadlineStr == "" {
				deadlineStr = "2026-05-24T23:59:59" // Fallback ke default deadline
			}
			if t, err := parseDeadline(deadlineStr); err == nil {
				if time.Now().After(t) {
					return entity.HackathonTeam{}, fmt.Errorf("batas waktu (deadline) pengumpulan proposal telah berakhir")
				}
			}
		}
		// Enforce Video Block & Deadline (Stage 2)
		if stage == "stage2" {
			if setting.HackathonVideoDisabled {
				return entity.HackathonTeam{}, fmt.Errorf("pengumpulan video (Stage 2) telah ditutup oleh administrator")
			}
			deadlineStr := setting.HackathonVideoDeadline
			if deadlineStr == "" {
				deadlineStr = "2026-06-06T23:59:59" // Fallback ke default deadline
			}
			if t, err := parseDeadline(deadlineStr); err == nil {
				if time.Now().After(t) {
					return entity.HackathonTeam{}, fmt.Errorf("batas waktu (deadline) pengumpulan video telah berakhir")
				}
			}
		}
		// Enforce Final Block & Deadline (Final Stage)
		if stage == "final" {
			if setting.HackathonFinalDisabled {
				return entity.HackathonTeam{}, fmt.Errorf("pengumpulan presentasi final telah ditutup oleh administrator")
			}
			deadlineStr := setting.HackathonFinalDeadline
			if deadlineStr == "" {
				deadlineStr = "2026-06-20T23:59:59" // Fallback ke default deadline
			}
			if t, err := parseDeadline(deadlineStr); err == nil {
				if time.Now().After(t) {
					return entity.HackathonTeam{}, fmt.Errorf("batas waktu (deadline) pengumpulan presentasi final telah berakhir")
				}
			}
		}
	}

	// Enforce qualified active stage progression block
	if stage == "stage1" && submission.Stage != "Stage-1" && submission.Stage != "Registered" {
		return entity.HackathonTeam{}, fmt.Errorf("akses ditutup: tim Anda saat ini berada pada tahap %s", submission.Stage)
	}
	if stage == "stage2" && submission.Stage != "Stage-2" {
		return entity.HackathonTeam{}, fmt.Errorf("akses ditutup: tim Anda saat ini berada pada tahap %s", submission.Stage)
	}
	if stage == "final" && submission.Stage != "Final" {
		return entity.HackathonTeam{}, fmt.Errorf("akses ditutup: tim Anda saat ini berada pada tahap %s", submission.Stage)
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

	// Preload the Team relation to ensure team details are populated for auditing
	if err := s.db.Preload("Team").First(&submission, submission.ID_HackathonTeam).Error; err != nil {
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
