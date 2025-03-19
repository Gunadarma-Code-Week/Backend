package service

import (
	"gcw/entity"
	"gcw/repository"
)

type UserService struct {
	userRepository repository.UserRepository
}

// type ProfileService interface {
// 	Get(uint64) (*entity.User, error)
// }

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		userRepository: repo,
	}
}

func (s *UserService) Update(user *entity.User, id uint64) error {
	user.ProfileHasUpdated = true
	return s.userRepository.Update(user, id)
}

func (s *UserService) FindById(id uint64) (*entity.User, error) {
	user := &entity.User{}
	err := s.userRepository.FindById(id, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
