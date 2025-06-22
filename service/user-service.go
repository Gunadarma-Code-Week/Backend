package service

import (
	"gcw/dto"
	"gcw/entity"
	"gcw/repository"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	userRepository repository.UserRepository
	DB             *gorm.DB
}

// type ProfileService interface {
// 	Get(uint64) (*entity.User, error)
// }

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		userRepository: repo,
		DB:             repo.GetDB(),
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

// getStatusOrUnregistered returns "Unregistered" if the status is empty, otherwise returns the status
func getStatusOrUnregistered(status string) string {
	if status == "" {
		return "Unregistered"
	}
	return status
}

func (s *UserService) GetEvents(userId uint64) (dto.ResponseEvents, error) {
	// Convert uint64 to string for the query
	idUserStr := strconv.FormatUint(userId, 10)

	var user entity.User
	if err := s.DB.Preload("Team").Where("id = ?", idUserStr).First(&user).Error; err != nil {
		return dto.ResponseEvents{}, err
	}

	// Get all team members
	var member []entity.User
	if err := s.DB.Where("id_team = ?", user.IDTeam).Find(&member).Error; err != nil {
		return dto.ResponseEvents{}, err
	}

	// Create user response
	responseUser := dto.User{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		University: user.Institusi,
	}

	// Create members response
	var responseMembers []dto.Member
	for _, m := range member {
		role := "Member"
		if m.ID == user.Team.ID_LeadTeam {
			role = "Leader"
		}
		responseMembers = append(responseMembers, dto.Member{
			Name:  m.Name,
			Role:  role,
			Email: m.Email,
		})
	}

	// Get event data
	var hackaton entity.HackathonTeam
	var cp entity.CPTeam
	var seminar entity.Seminar

	_ = s.DB.Where("id_team = ?", user.IDTeam).First(&hackaton)
	_ = s.DB.Where("id_team = ?", user.IDTeam).First(&cp)
	_ = s.DB.Where("id_user = ?", idUserStr).First(&seminar)

	seminarStatus := "Unregistered"
	if seminar.ID_Seminar != 0 {
		seminarStatus = "Registered"
	}

	// Fill event data
	events := []dto.Event{
		{
			Name:   "Seminar",
			Status: seminarStatus,
			Ticket: dto.Ticket{},
		},
		{
			Name:   "Hackathon",
			Status: getStatusOrUnregistered(hackaton.Stage),
			Ticket: dto.Ticket{},
		},
		{
			Name:   "Competitive Programming",
			Status: getStatusOrUnregistered(cp.Stage),
			Ticket: dto.Ticket{},
		},
	}

	return dto.ResponseEvents{
		User:   responseUser,
		Events: events,
		IdTeam: user.Team.JoinCode,
	}, nil
}
