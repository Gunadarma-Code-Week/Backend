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
	Create(*dto.RegistrationResponseDTO) (*entity.Team, error)
}

func GatRegistrationService(repo repository.RegistrationRepository) RegistrationService {
	return &registrationService{
		registrationRepository: repo,
	}
}

func (s *registrationService) Create(registrationDTO *dto.RegistrationResponseDTO) (*entity.Team, error) {
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
