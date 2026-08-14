package service

import (
	"errors"
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/repository"

	"strings"
	"time"

	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

type RegistrationService struct {
	registrationRepository *repository.RegistrationRepository
	userRepository         *repository.UserRepository
	domJudgeService        *DomJudgeService
}

func NewRegistrationService(
	rp *repository.RegistrationRepository,
	ds *DomJudgeService,
) *RegistrationService {
	return &RegistrationService{
		registrationRepository: rp,
		domJudgeService:        ds,
	}
}

func (s *RegistrationService) Repository() *repository.RegistrationRepository {
	return s.registrationRepository
}

func (s *RegistrationService) CPTeamRegistration(
	registrationDTO *dto.RegistrationCPTeamRequest,
	userLead *entity.User,
) (*dto.RegistrationCPTeamResponse, error) {
	var err error
	registrationDTO.TeamName = strings.TrimSpace(registrationDTO.TeamName)
	teamRegistration := &entity.Team{
		TeamName:       registrationDTO.TeamName,
		Supervisor:     registrationDTO.Supervisor,
		SupervisorNIDN: registrationDTO.SupervisorNIDN,
		ID_LeadTeam:    userLead.ID,
		Event:          "cp",
	}

	// CP Fee: Rp 50,000
	const cpFee = 50000

	if userLead.IDTeam != nil {
		logging.Low("RegistrationService.CPTeamRegistration", "BAD_REQUEST", "User already have team")
		return nil, fmt.Errorf("USER ALREADY HAVE TEAM")
	}

	// Check if CP registration is disabled by admin
	var setting entity.SystemSetting
	if err = s.registrationRepository.DB.First(&setting, 1).Error; err == nil {
		if setting.CPRegistrationDisabled {
			logging.Low("RegistrationService.CPTeamRegistration", "BAD_REQUEST", "Pendaftaran Competitive Programming telah ditutup oleh administrator")
			return nil, fmt.Errorf("Pendaftaran Competitive Programming telah ditutup")
		}
		if setting.CPRegistrationDeadline != nil && *setting.CPRegistrationDeadline != "" {
			if parsedTime, parseErr := time.Parse(time.RFC3339, *setting.CPRegistrationDeadline); parseErr == nil {
				if time.Now().After(parsedTime) {
					logging.Low("RegistrationService.CPTeamRegistration", "BAD_REQUEST", "Batas waktu pendaftaran Competitive Programming telah berakhir")
					return nil, fmt.Errorf("Batas waktu pendaftaran Competitive Programming telah berakhir")
				}
			}
		}
	}

	// Check duplicate team name secara global (hackathon, cp, ctf)
	if err = s.registrationRepository.FindActiveTeamByNameGlobal(registrationDTO.TeamName); err == nil {
		logging.Low("RegistrationService.CPTeamRegistration", "BAD_REQUEST", "Nama Tim Sudah Digunakan")
		return nil, fmt.Errorf("Nama Tim Sudah Digunakan")
	}

	// generate join code
	var joinCode string

	for {
		joinCode = helper.RandomStringNumber(6)
		err = s.registrationRepository.FindTeamByJoinCode(&entity.Team{}, joinCode)
		if err != nil {
			break
		}
	}

	teamRegistration.JoinCode = joinCode
	// Manual payment logic
	orderID := fmt.Sprintf("CP-MANUAL-%d-%d", userLead.ID, time.Now().UnixNano())
	teamRegistration.KomitmenFee = registrationDTO.BuktiPembayaran
	qrString := "-"

	teamRegistration.OrderID = orderID
	teamRegistration.QRString = qrString

	// create domjudge team first before creating team
	domJudgeUsername, domJudgePassword, err := s.domJudgeService.CreateDomJudgeTeamUser(
		joinCode,
		registrationDTO.TeamName,
		userLead.Email,
	)
	if err != nil && domJudgeUsername != "skipped" {
		logging.Low("RegistrationService.CPTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	tx := s.registrationRepository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// create team
	err = s.registrationRepository.CreateTeam(tx, teamRegistration)
	if err != nil {
		logging.Low("RegistrationService.CPTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		tx.Rollback()
		return nil, err
	}

	// create cp team
	cpTeam := &entity.CPTeam{
		IDTeam:           teamRegistration.ID_Team,
		Stage:            "Registered",
		Status:           "Registration",
		DomjudgeUsername: domJudgeUsername,
		DomjudgePassword: domJudgePassword,
	}
	err = s.registrationRepository.CreateCPTeam(tx, cpTeam)
	if err != nil {
		logging.Low("RegistrationService.CPTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		tx.Rollback()
		return nil, err
	}

	// update user team id
	err = s.registrationRepository.UpdateUserTeam(tx, userLead, teamRegistration.ID_Team, userLead.ID)
	if err != nil {
		logging.Low("RegistrationService.CPTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		logging.Low("RegistrationService.CPTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	registrationTeamResponse := &dto.RegistraionTeamResponse{}
	err = smapping.FillStruct(registrationTeamResponse, smapping.MapFields(teamRegistration))
	if err != nil {
		logging.Low("RegistrationService.CPTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	registrasionCPResponse := &dto.RegistrationCPResponse{}
	registrasionCPResponse.JoinCode = joinCode
	registrasionCPResponse.QRString = qrString
	registrasionCPResponse.OrderID = teamRegistration.OrderID
	registrasionCPResponse.PaymentStatus = teamRegistration.PaymentStatus
	err = smapping.FillStruct(registrasionCPResponse, smapping.MapFields(cpTeam))
	if err != nil {
		logging.Low("RegistrationService.CPTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	registrationCPTeamResponse := &dto.RegistrationCPTeamResponse{
		Team:   *registrationTeamResponse,
		CPTeam: *registrasionCPResponse,
	}

	return registrationCPTeamResponse, nil
}

func (s *RegistrationService) HackathonTeamRegistration(
	registrationDTO *dto.RegistrationHackathonTeamRequest,
	userLead *entity.User,
) (*dto.RegistrationHackathonTeamResponse, error) {
	var err error
	registrationDTO.TeamName = strings.TrimSpace(registrationDTO.TeamName)
	if userLead.IDTeam != nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "BAD_REQUEST", "User already have team")
		return nil, fmt.Errorf("USER ALREADY HAVE TEAM")
	}

	// Check if Hackathon registration is disabled by admin
	var setting entity.SystemSetting
	if err = s.registrationRepository.DB.First(&setting, 1).Error; err == nil {
		if setting.HackathonRegistrationDisabled {
			logging.Low("RegistrationService.HackathonTeamRegistration", "BAD_REQUEST", "Pendaftaran Hackathon telah ditutup oleh administrator")
			return nil, fmt.Errorf("Pendaftaran Hackathon telah ditutup")
		}
	}

	// Check duplicate team name secara global (hackathon, cp, ctf)
	if err = s.registrationRepository.FindActiveTeamByNameGlobal(registrationDTO.TeamName); err == nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "BAD_REQUEST", "Nama Tim Sudah Digunakan")
		return nil, fmt.Errorf("Nama Tim Sudah Digunakan")
	}

	teamRegistration := &entity.Team{
		TeamName:       registrationDTO.TeamName,
		Supervisor:     registrationDTO.Supervisor,
		SupervisorNIDN: registrationDTO.SupervisorNIDN,
		ID_LeadTeam:    userLead.ID,
		Event:          "hackathon",
	}

	// Hackathon Fee: Rp 120,000 (Normal Fee)
	const hackathonFee = 120000

	// generate join code
	var joinCode string

	for {
		joinCode = helper.RandomStringNumber(6)
		err = s.registrationRepository.FindTeamByJoinCode(&entity.Team{}, joinCode)
		if err != nil {
			break
		}
	}
	teamRegistration.JoinCode = joinCode

	// Manual payment logic
	orderID := fmt.Sprintf("HACK-MANUAL-%d-%d", userLead.ID, time.Now().UnixNano())
	teamRegistration.KomitmenFee = registrationDTO.BuktiPembayaran
	qrString := "-"

	teamRegistration.OrderID = orderID
	teamRegistration.QRString = qrString

	tx := s.registrationRepository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// create team
	err = s.registrationRepository.CreateTeam(tx, teamRegistration)
	if err != nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		tx.Rollback()
		return nil, err
	}

	// create hackathon team
	hackathonTeam := &entity.HackathonTeam{
		IDTeam: teamRegistration.ID_Team,
		Stage:  "Registered",
		Status: "Registration",
	}
	err = s.registrationRepository.CreateHackathonTeam(tx, hackathonTeam)
	if err != nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		tx.Rollback()
		return nil, err
	}

	// update user team id
	err = s.registrationRepository.UpdateUserTeam(tx, userLead, teamRegistration.ID_Team, userLead.ID)
	if err != nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	registrationTeamResponse := &dto.RegistraionTeamResponse{}
	err = smapping.FillStruct(registrationTeamResponse, smapping.MapFields(teamRegistration))
	if err != nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	registrasionHackathonResponse := &dto.RegistrationHackathonResponse{}
	registrasionHackathonResponse.JoinCode = joinCode
	registrasionHackathonResponse.QRString = qrString
	registrasionHackathonResponse.OrderID = orderID
	registrasionHackathonResponse.PaymentStatus = teamRegistration.PaymentStatus
	err = smapping.FillStruct(registrasionHackathonResponse, smapping.MapFields(hackathonTeam))
	if err != nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	registrationHackathonTeamResponse := &dto.RegistrationHackathonTeamResponse{
		Team:          *registrationTeamResponse,
		HackathonTeam: *registrasionHackathonResponse,
	}

	return registrationHackathonTeamResponse, nil
}

func (s *RegistrationService) CTFTeamRegistration(
	registrationDTO *dto.RegistrationCTFTeamRequest,
	userLead *entity.User,
) (*dto.RegistrationCTFTeamResponse, error) {
	var err error
	registrationDTO.TeamName = strings.TrimSpace(registrationDTO.TeamName)
	teamRegistration := &entity.Team{
		TeamName:       registrationDTO.TeamName,
		Supervisor:     registrationDTO.Supervisor,
		SupervisorNIDN: registrationDTO.SupervisorNIDN,
		ID_LeadTeam:    userLead.ID,
		Event:          "ctf",
	}

	// CTF Fee: Rp 75,000
	// const ctfFee = 75000

	if userLead.IDTeam != nil {
		logging.Low("RegistrationService.CTFTeamRegistration", "BAD_REQUEST", "User already have team")
		return nil, fmt.Errorf("USER ALREADY HAVE TEAM")
	}

	// Check if CTF registration is disabled by admin
	var setting entity.SystemSetting
	if err = s.registrationRepository.DB.First(&setting, 1).Error; err == nil {
		if setting.CTFRegistrationDisabled {
			logging.Low("RegistrationService.CTFTeamRegistration", "BAD_REQUEST", "Pendaftaran Capture The Flag telah ditutup oleh administrator")
			return nil, fmt.Errorf("Pendaftaran Capture The Flag telah ditutup")
		}
		if setting.CTFRegistrationDeadline != nil && *setting.CTFRegistrationDeadline != "" {
			if parsedTime, parseErr := time.Parse(time.RFC3339, *setting.CTFRegistrationDeadline); parseErr == nil {
				if time.Now().After(parsedTime) {
					logging.Low("RegistrationService.CTFTeamRegistration", "BAD_REQUEST", "Batas waktu pendaftaran Capture The Flag telah berakhir")
					return nil, fmt.Errorf("Batas waktu pendaftaran Capture The Flag telah berakhir")
				}
			}
		}
	}

	// Check duplicate team name secara global (hackathon, cp, ctf)
	if err = s.registrationRepository.FindActiveTeamByNameGlobal(registrationDTO.TeamName); err == nil {
		logging.Low("RegistrationService.CTFTeamRegistration", "BAD_REQUEST", "Nama Tim Sudah Digunakan")
		return nil, fmt.Errorf("Nama Tim Sudah Digunakan")
	}

	// generate join code
	var joinCode string

	for {
		joinCode = helper.RandomStringNumber(6)
		err = s.registrationRepository.FindTeamByJoinCode(&entity.Team{}, joinCode)
		if err != nil {
			break
		}
	}

	teamRegistration.JoinCode = joinCode
	teamRegistration.KomitmenFee = registrationDTO.BuktiPembayaran

	// Manual payment logic
	orderID := fmt.Sprintf("CTF-MANUAL-%d-%d", userLead.ID, time.Now().UnixNano())
	teamRegistration.KomitmenFee = registrationDTO.BuktiPembayaran
	qrString := "-"

	teamRegistration.OrderID = orderID
	teamRegistration.QRString = qrString

	// create domjudge team first before creating team
	domJudgeUsername, domJudgePassword, err := s.domJudgeService.CreateDomJudgeTeamUser(
		joinCode,
		registrationDTO.TeamName,
		userLead.Email,
	)
	if err != nil && domJudgeUsername != "skipped" {
		logging.Low("RegistrationService.CTFTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	tx := s.registrationRepository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// create team
	err = s.registrationRepository.CreateTeam(tx, teamRegistration)
	if err != nil {
		logging.Low("RegistrationService.CTFTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		tx.Rollback()
		return nil, err
	}

	// create ctf team
	ctfTeam := &entity.CTFTeam{
		IDTeam:           teamRegistration.ID_Team,
		Stage:            "Registered",
		Status:           "Registration",
		DomjudgeUsername: domJudgeUsername,
		DomjudgePassword: domJudgePassword,
	}
	err = s.registrationRepository.CreateCTFTeam(tx, ctfTeam)
	if err != nil {
		logging.Low("RegistrationService.CTFTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		tx.Rollback()
		return nil, err
	}

	// update user team id
	err = s.registrationRepository.UpdateUserTeam(tx, userLead, teamRegistration.ID_Team, userLead.ID)
	if err != nil {
		logging.Low("RegistrationService.CTFTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		logging.Low("RegistrationService.CTFTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	registrationTeamResponse := &dto.RegistraionTeamResponse{}
	err = smapping.FillStruct(registrationTeamResponse, smapping.MapFields(teamRegistration))
	if err != nil {
		logging.Low("RegistrationService.CTFTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	registrasionCTFResponse := &dto.RegistrationCTFResponse{}
	registrasionCTFResponse.JoinCode = joinCode
	registrasionCTFResponse.QRString = qrString
	registrasionCTFResponse.OrderID = teamRegistration.OrderID
	registrasionCTFResponse.PaymentStatus = teamRegistration.PaymentStatus
	err = smapping.FillStruct(registrasionCTFResponse, smapping.MapFields(ctfTeam))
	if err != nil {
		logging.Low("RegistrationService.CTFTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	registrationCTFTeamResponse := &dto.RegistrationCTFTeamResponse{
		Team:    *registrationTeamResponse,
		CTFTeam: *registrasionCTFResponse,
	}

	return registrationCTFTeamResponse, nil
}

func (s *RegistrationService) FindTeamByJoinCode(joinCode string) (*entity.Team, error) {
	team := &entity.Team{}
	err := s.registrationRepository.FindTeamByJoinCode(team, joinCode)
	if err != nil {
		logging.Low("RegistrationService.FindTeamByJoinCode", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	return team, nil
}

func (s *RegistrationService) JoinTeam(
	joinCode string,
	user *entity.User,
) (*entity.Team, error) {
	team := &entity.Team{}
	err := s.registrationRepository.FindTeamByJoinCode(team, joinCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logging.Low("RegistrationService.JoinTeam", "NOT_FOUND", "Team not found")
			return nil, fmt.Errorf("TEAM NOT FOUND")
		}
		logging.Low("RegistrationService.JoinTeam", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	// Check if registration is disabled by admin for this team's event
	var setting entity.SystemSetting
	if err = s.registrationRepository.DB.First(&setting, 1).Error; err == nil {
		if team.Event == "hackathon" && setting.HackathonRegistrationDisabled {
			logging.Low("RegistrationService.JoinTeam", "BAD_REQUEST", "Pendaftaran Hackathon telah ditutup oleh administrator")
			return nil, fmt.Errorf("Pendaftaran Hackathon telah ditutup")
		}
		if team.Event == "cp" && setting.CPRegistrationDisabled {
			logging.Low("RegistrationService.JoinTeam", "BAD_REQUEST", "Pendaftaran Competitive Programming telah ditutup oleh administrator")
			return nil, fmt.Errorf("Pendaftaran Competitive Programming telah ditutup")
		}
		if team.Event == "ctf" && setting.CTFRegistrationDisabled {
			logging.Low("RegistrationService.JoinTeam", "BAD_REQUEST", "Pendaftaran Capture The Flag telah ditutup oleh administrator")
			return nil, fmt.Errorf("Pendaftaran Capture The Flag telah ditutup")
		}
		if team.Event == "cp" && setting.CPRegistrationDeadline != nil && *setting.CPRegistrationDeadline != "" {
			if parsedTime, parseErr := time.Parse(time.RFC3339, *setting.CPRegistrationDeadline); parseErr == nil {
				if time.Now().After(parsedTime) {
					logging.Low("RegistrationService.JoinTeam", "BAD_REQUEST", "Batas waktu pendaftaran Competitive Programming telah berakhir")
					return nil, fmt.Errorf("Batas waktu pendaftaran Competitive Programming telah berakhir")
				}
			}
		}
		if team.Event == "ctf" && setting.CTFRegistrationDeadline != nil && *setting.CTFRegistrationDeadline != "" {
			if parsedTime, parseErr := time.Parse(time.RFC3339, *setting.CTFRegistrationDeadline); parseErr == nil {
				if time.Now().After(parsedTime) {
					logging.Low("RegistrationService.JoinTeam", "BAD_REQUEST", "Batas waktu pendaftaran Capture The Flag telah berakhir")
					return nil, fmt.Errorf("Batas waktu pendaftaran Capture The Flag telah berakhir")
				}
			}
		}
	}

	userCount, err := s.registrationRepository.CountUserByTeamID(team.ID_Team)
	if err != nil {
		logging.Low("RegistrationService.JoinTeam", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	isTeamFul :=
		team.Event == "hackathon" && userCount >= 5 ||
			(team.Event == "cp" || team.Event == "ctf") && userCount >= 3
	if isTeamFul {
		logging.Low("RegistrationService.JoinTeam", "BAD_REQUEST", "Team is full")
		return nil, fmt.Errorf("TEAM IS FULL")
	}

	if user.IDTeam != nil {
		logging.Low("RegistrationService.JoinTeam", "BAD_REQUEST", "User already have team")
		return nil, fmt.Errorf("USER ALREADY HAVE TEAM")
	}

	// NOTE : better create new function for this instead use tx func
	db := s.registrationRepository.DB
	err = s.registrationRepository.UpdateUserTeam(db, user, team.ID_Team, user.ID)
	if err != nil {
		logging.Low("RegistrationService.JoinTeam", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	return team, nil
}

func (s *RegistrationService) UpdatePaymentStatus(orderID string, status string) error {
	team := &entity.Team{}
	err := s.registrationRepository.DB.Where("order_id = ?", orderID).First(team).Error
	if err != nil {
		return err
	}

	team.PaymentStatus = status

	tx := s.registrationRepository.DB.Begin()
	if err := tx.Save(team).Error; err != nil {
		tx.Rollback()
		return err
	}

	// If paid, update the event-specific team status to Verified and stage to Stage-1
	if status == "Paid" {
		if team.Event == "hackathon" {
			err = tx.Model(&entity.HackathonTeam{}).Where("id_team = ?", team.ID_Team).Updates(map[string]interface{}{
				"status": "Verified",
				"stage":  "Stage-1",
			}).Error
		} else if team.Event == "cp" {
			err = tx.Model(&entity.CPTeam{}).Where("id_team = ?", team.ID_Team).Updates(map[string]interface{}{
				"status": "Verified",
				"stage":  "Stage-1",
			}).Error
		} else if team.Event == "ctf" {
			err = tx.Model(&entity.CTFTeam{}).Where("id_team = ?", team.ID_Team).Updates(map[string]interface{}{
				"status": "Verified",
				"stage":  "Stage-1",
			}).Error
		}
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *RegistrationService) UpdateReceipt(orderID string, receiptURL string) error {
	team := &entity.Team{}
	err := s.registrationRepository.DB.Where("order_id = ?", orderID).First(team).Error
	if err != nil {
		return err
	}

	team.KomitmenFee = receiptURL
	return s.registrationRepository.DB.Save(team).Error
}

func (s *RegistrationService) UpdateTeamDetails(orderID string, teamName string, supervisor string, supervisorNIDN string, receiptURL string) (string, error) {
	teamName = strings.TrimSpace(teamName)
	team := &entity.Team{}
	err := s.registrationRepository.DB.Where("order_id = ?", orderID).First(team).Error
	if err != nil {
		return "", err
	}

	var changes []string

	// Validation: Check if the new team name is already taken by ANOTHER team
	if !strings.EqualFold(team.TeamName, teamName) {
		err = s.registrationRepository.FindActiveTeamByNameGlobalExcludingTeam(teamName, team.ID_Team)
		fmt.Printf("[DEBUG] Checking team name uniqueness for '%s' (excluding ID %d). Result error: %v\n", teamName, team.ID_Team, err)
		if err == nil {
			logging.Low("RegistrationService.UpdateTeamDetails", "BAD_REQUEST", "Nama Tim Sudah Digunakan")
			return "", fmt.Errorf("Nama Tim Sudah Digunakan")
		}
		changes = append(changes, fmt.Sprintf("Nama Tim ('%s' -> '%s')", team.TeamName, teamName))
	}

	if team.Supervisor != supervisor {
		changes = append(changes, fmt.Sprintf("Pembimbing ('%s' -> '%s')", team.Supervisor, supervisor))
	}
	if team.SupervisorNIDN != supervisorNIDN {
		changes = append(changes, fmt.Sprintf("NIDN ('%s' -> '%s')", team.SupervisorNIDN, supervisorNIDN))
	}
	if team.KomitmenFee != receiptURL {
		changes = append(changes, fmt.Sprintf("Bukti Pembayaran ('%s' -> '%s')", team.KomitmenFee, receiptURL))
	}

	team.TeamName = teamName
	team.Supervisor = supervisor
	team.SupervisorNIDN = supervisorNIDN
	team.KomitmenFee = receiptURL

	err = s.registrationRepository.DB.Save(team).Error
	if err != nil {
		return "", err
	}

	return strings.Join(changes, ", "), nil
}
