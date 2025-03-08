package service

import (
	"gcw/dto"
	"gcw/entity"
	"gcw/helper/logging"
	"gcw/repository"

	"github.com/mashingan/smapping"
)

type profileService struct {
	profileRepository repository.ProfileRepository
}

type ProfileService interface {
	Get(uint64) (*entity.Profile, error)
	Create(*dto.ProfileResponseDTO) (*entity.Profile, error)
}

func GateProfileService(repo repository.ProfileRepository) ProfileService {
	return &profileService{
		profileRepository: repo,
	}
}

func (s *profileService) Get(id uint64) (*entity.Profile, error) {
	profile := &entity.Profile{}

	err := s.profileRepository.FindById(id, profile)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *profileService) Create(profileDTO *dto.ProfileResponseDTO) (*entity.Profile, error) {
	profile := &entity.Profile{}

	if err := smapping.FillStruct(profile, smapping.MapFields(profileDTO)); err != nil {
		logging.Low("ProfileService.Create", "BAD_REQUEST", err.Error())
		return nil, err
	}
	if err := s.profileRepository.Create(profile); err != nil {
		logging.Warn("ProfileService.Create", "INTERNAL_SERVER_ERROR", err.Error())
		return nil, err
	}
	return profile, nil
}
