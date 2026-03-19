package service

import (
	"errors"
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/repository"

	"github.com/mashingan/smapping"
	"gorm.io/gorm"
	"time"
)

type RegistrationService struct {
	registrationRepository *repository.RegistrationRepository
	userRepository         *repository.UserRepository
	domJudgeService        *DomJudgeService
	midtransService        *MidtransService
}

func NewRegistrationService(
	rp *repository.RegistrationRepository,
	ds *DomJudgeService,
	ms *MidtransService,
) *RegistrationService {
	return &RegistrationService{
		registrationRepository: rp,
		domJudgeService:        ds,
		midtransService:        ms,
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

	// Check duplicate team name
	if err = s.registrationRepository.FindTeamByNameAndEvent(&entity.Team{}, registrationDTO.TeamName, "cp"); err == nil {
		logging.Low("RegistrationService.CPTeamRegistration", "BAD_REQUEST", "Team name already taken")
		return nil, fmt.Errorf("TEAM NAME ALREADY TAKEN")
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

	// Generate Order ID for CP (manual payment)
	orderID := fmt.Sprintf("CP-%d-%d", userLead.ID, time.Now().UnixNano())
	teamRegistration.OrderID = orderID
	teamRegistration.QRString = "-"

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
	registrasionCPResponse.QRString = "-"
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
	if userLead.IDTeam != nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "BAD_REQUEST", "User already have team")
		return nil, fmt.Errorf("USER ALREADY HAVE TEAM")
	}

	// Check duplicate team name
	if err = s.registrationRepository.FindTeamByNameAndEvent(&entity.Team{}, registrationDTO.TeamName, "hackathon"); err == nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "BAD_REQUEST", "Team name already taken")
		return nil, fmt.Errorf("TEAM NAME ALREADY TAKEN")
	}

	teamRegistration := &entity.Team{
		TeamName:       registrationDTO.TeamName,
		Supervisor:     registrationDTO.Supervisor,
		SupervisorNIDN: registrationDTO.SupervisorNIDN,
		ID_LeadTeam:    userLead.ID,
		Event:          "hackathon",
	}

	// Hackathon Fee: Rp 100,000
	const hackathonFee = 100000

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

	// Generate Midtrans Order ID and QRIS
	orderID := fmt.Sprintf("HACK-%d-%d", userLead.ID, time.Now().UnixNano())
	teamRegistration.OrderID = orderID

	qrString := "-"
	if s.midtransService != nil {
		qrString, err = s.midtransService.GenerateQRIS(orderID, int64(hackathonFee))
		if err != nil {
			logging.Low("RegistrationService.HackathonTeamRegistration", "INTERNAL_SERVER_ERROR", "Midtrans QR Generation failed: "+err.Error())
			return nil, err
		}
	}
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

	// Check duplicate team name
	if err = s.registrationRepository.FindTeamByNameAndEvent(&entity.Team{}, registrationDTO.TeamName, "ctf"); err == nil {
		logging.Low("RegistrationService.CTFTeamRegistration", "BAD_REQUEST", "Team name already taken")
		return nil, fmt.Errorf("TEAM NAME ALREADY TAKEN")
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

	// Generate Order ID for CTF (manual payment)
	orderID := fmt.Sprintf("CTF-%d-%d", userLead.ID, time.Now().UnixNano())
	teamRegistration.OrderID = orderID
	teamRegistration.QRString = "-"

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
		IDTeam: teamRegistration.ID_Team,
		Stage:  "Registered",
		Status: "Registration",
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
	registrasionCTFResponse.QRString = "-"
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

