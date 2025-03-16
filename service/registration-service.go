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
