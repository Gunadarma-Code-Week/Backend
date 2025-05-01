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

func (s *RegistrationService) CPTeamRegistration(
	registrationDTO *dto.RegistrationCPTeamRequest,
	userLead *entity.User,
) (*dto.RegistrationCPTeamResponse, error) {
	teamRegistration := &entity.Team{
		TeamName:       registrationDTO.TeamName,
		Supervisor:     registrationDTO.Supervisor,
		SupervisorNIDN: registrationDTO.SupervisorNIDN,
		ID_LeadTeam:    userLead.ID,
		KomitmenFee:    registrationDTO.KomitmenFee,
		Event:          "cp",
	}

	if userLead.IDTeam != nil {
		logging.Low("RegistrationService.CPTeamRegistration", "BAD_REQUEST", "User already have team")
		return nil, fmt.Errorf("USER ALREADY HAVE TEAM")
	}

	var joinCode string

	for {
		joinCode = helper.RandomStringNumber(6)
		err := s.registrationRepository.FindTeamByJoinCode(&entity.Team{}, joinCode)
		if err != nil {
			break
		}
	}
	teamRegistration.JoinCode = joinCode

	// create domjudge team first before creating team
	domJudgeUsername, domJudgePassword, err := s.domJudgeService.CreateDomJudgeTeamUser(
		joinCode,
		registrationDTO.TeamName,
		userLead.Email,
	)
	if err != nil {
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
		Stage:            "Registration",
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
	registrasionCPResponse.KomitmenFee = teamRegistration.KomitmenFee
	registrasionCPResponse.JoinCode = joinCode
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
	if userLead.IDTeam != nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "BAD_REQUEST", "User already have team")
		return nil, fmt.Errorf("USER ALREADY HAVE TEAM")
	}

	teamRegistration := &entity.Team{
		TeamName:       registrationDTO.TeamName,
		Supervisor:     registrationDTO.Supervisor,
		SupervisorNIDN: registrationDTO.SupervisorNIDN,
		ID_LeadTeam:    userLead.ID,
		KomitmenFee:    registrationDTO.KomitmenFee,
		Event:          "hackathon",
	}

	// generate join code
	var joinCode string

	for {
		joinCode = helper.RandomStringNumber(6)
		err := s.registrationRepository.FindTeamByJoinCode(&entity.Team{}, joinCode)
		if err != nil {
			break
		}
	}
	teamRegistration.JoinCode = joinCode

	tx := s.registrationRepository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// create team
	err := s.registrationRepository.CreateTeam(tx, teamRegistration)
	if err != nil {
		logging.Low("RegistrationService.HackathonTeamRegistration", "INTERNAL_SERVER_ERROR", err.Error())
		tx.Rollback()
		return nil, err
	}

	// create hackathon team
	hackathonTeam := &entity.HackathonTeam{
		IDTeam: teamRegistration.ID_Team,
		Stage:  "Registration",
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
	registrasionHackathonResponse.KomitmenFee = teamRegistration.KomitmenFee
	registrasionHackathonResponse.JoinCode = joinCode
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
			team.Event == "cp" && userCount >= 3
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
