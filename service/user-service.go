package service

import (
	"gcw/dto"
	"gcw/entity"
	"gcw/repository"
	"time"
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

func (s *UserService) FindByIdTeam(id, id_leader uint64) ([]dto.Member, error) {
	var teamMembers []entity.User
	err := s.userRepository.FindByIdTeam(id, &teamMembers)
	if err != nil {
		return nil, err
	}

	members := []dto.Member{}
	for _, data := range teamMembers {
		if data.ID == id_leader {
			continue
		}

		member := dto.Member{}

		member.Name = data.Name
		member.Email = data.Email
		member.Role = data.Role

		members = append(members, member)
	}

	return members, nil
}

func (s *UserService) GetUsersByDateRange(startDate, endDate time.Time, limit, offset int) ([]*entity.User, int64, error) {
	// Fetch users from the repository with date range and pagination
	users, totalUsers, err := s.userRepository.FindAll(startDate, endDate, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, totalUsers, nil
}
