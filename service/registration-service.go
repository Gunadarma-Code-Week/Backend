package service

import (
	"gcw/dto"
	"gcw/entity"
	"gcw/helper/logging"
	"gcw/repository"

	"github.com/mashingan/smapping"
)

type registrationService struct {
	registrationRepository repository.RegistrationRepository
}

type RegistrationService interface {
	Create(*dto.RegistrationResponseWithJoinCode) (*entity.Team, error)
	CreateTeam(*dto.RegistrationResponseHackathon) (*entity.HackathonTeam, error)
	UpdateUser(uint64, uint64) (*entity.User, error)
	UpdateUserJoinCode(string, uint64) (*entity.User, error)
}

func GatRegistrationService(repo repository.RegistrationRepository) RegistrationService {
	return &registrationService{
		registrationRepository: repo,
	}
}

func (s *registrationService) Create(registrationDTO *dto.RegistrationResponseWithJoinCode) (*entity.Team, error) {
	registration := &entity.Team{}

	if err := smapping.FillStruct(registration, smapping.MapFields(registrationDTO)); err != nil {
		logging.Low("RegistrationService.Create", "BAD_REQUEST", err.Error())
		return nil, err
	}

	if err := s.registrationRepository.Create(registration); err != nil {
		logging.Warn("RegistrationService.Create", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	return registration, nil
}

func (s *registrationService) UpdateUser(id_team uint64, id_user uint64) (*entity.User, error) {
	user := &entity.User{}

	if err := s.registrationRepository.UpdateUserTeam(user, id_team, id_user); err != nil {
		logging.Warn("RegistrationService.Create", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	return user, nil
}

func (s *registrationService) UpdateUserJoinCode(code string, id_user uint64) (*entity.User, error) {
	user := &entity.User{}

	if err := s.registrationRepository.UpdateUserTeamJoinCode(user, code, id_user); err != nil {
		logging.Warn("RegistrationService.Create", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}

	return user, nil
}

func (s *registrationService) CreateTeam(registrationDTO *dto.RegistrationResponseHackathon) (*entity.HackathonTeam, error) {
	registration := &entity.HackathonTeam{}

	if err := smapping.FillStruct(registration, smapping.MapFields(registrationDTO)); err != nil {
		logging.Low("RegistrationService.Create", "BAD_REQUEST", err.Error())
		return nil, err
	}

	if err := s.registrationRepository.CreateTeam(registration); err != nil {
		logging.Warn("RegistrationService.Create", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}
	return registration, nil
}
