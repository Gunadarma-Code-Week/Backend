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

// Admin User Management Methods

// AdminGetAllUsers - Get all users with pagination, filtering, and sorting for admin
func (s *UserService) AdminGetAllUsers(query dto.AdminGetUsersQueryDTO) (dto.AdminUsersListResponseDTO, error) {
	var users []entity.User
	var totalUsers int64

	// Build query
	db := s.DB.Model(&entity.User{})

	// Apply date filters if provided
	if query.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", query.StartDate)
		if err == nil {
			db = db.Where("created_at >= ?", startDate)
		}
	}
	if query.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", query.EndDate)
		if err == nil {
			db = db.Where("created_at <= ?", endDate.Add(24*time.Hour))
		}
	}

	// Apply search filter if provided
	if query.Q != "" {
		searchTerm := "%" + query.Q + "%"
		db = db.Where("name ILIKE ? OR email ILIKE ? OR institusi ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Count total users
	if err := db.Count(&totalUsers).Error; err != nil {
		return dto.AdminUsersListResponseDTO{}, err
	}

	// Apply sorting
	sortBy := "id"
	if query.SortBy != "" {
		validSortFields := map[string]bool{
			"id": true, "institusi": true, "id_team": true, "nim": true,
			"soc_med_document": true, "profile_has_updated": true, "data_has_verified": true,
		}
		if validSortFields[query.SortBy] {
			sortBy = query.SortBy
		}
	}

	sortOrder := "ASC"
	if query.SortOrder == "DESC" {
		sortOrder = "DESC"
	}

	// Apply pagination and sorting
	offset := (query.Page - 1) * query.Limit
	if err := db.Order(sortBy + " " + sortOrder).Limit(query.Limit).Offset(offset).Find(&users).Error; err != nil {
		return dto.AdminUsersListResponseDTO{}, err
	}

	// Convert to response DTOs
	var userResponses []dto.AdminUserResponseDTO
	for _, user := range users {
		userResponse := dto.AdminUserResponseDTO{
			ID:                user.ID,
			Email:             user.Email,
			Role:              user.Role,
			Name:              user.Name,
			Institusi:         user.Institusi,
			Phone:             user.Phone,
			Jenjang:           user.Jenjang,
			Major:             user.Major,
			ProfilePicture:    user.ProfilePicture,
			ProfileHasUpdated: user.ProfileHasUpdated,
			DataHasVerified:   user.DataHasVerified,
			CreatedAt:         user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:         user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Handle nullable fields
		if user.NIM != nil {
			userResponse.NIM = *user.NIM
		}
		if user.Gender != nil {
			userResponse.Gender = *user.Gender
		}
		if user.BirthPlace != nil {
			userResponse.BirthPlace = *user.BirthPlace
		}
		if user.BirthDate != nil {
			userResponse.BirthDate = user.BirthDate.Format("2006-01-02")
		}
		if user.IDTeam != nil {
			userResponse.IDTeam = *user.IDTeam
		}
		userResponse.SocMedDocument = user.SocMedDocument
		userResponse.DokumenFilename = user.DokumenFilename

		userResponses = append(userResponses, userResponse)
	}

	// Calculate pagination info
	totalPages := (totalUsers + int64(query.Limit) - 1) / int64(query.Limit)
	hasMore := int64(query.Page) < totalPages

	return dto.AdminUsersListResponseDTO{
		Users: userResponses,
		Meta: dto.AdminUsersMetaDTO{
			TotalItems:  totalUsers,
			TotalPages:  totalPages,
			CurrentPage: query.Page,
			Limit:       query.Limit,
			HasMore:     hasMore,
		},
	}, nil
}

// AdminGetUserById - Get user by ID for admin
func (s *UserService) AdminGetUserById(id uint64) (*dto.AdminUserResponseDTO, error) {
	var user entity.User
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	userResponse := &dto.AdminUserResponseDTO{
		ID:                user.ID,
		Email:             user.Email,
		Role:              user.Role,
		Name:              user.Name,
		Institusi:         user.Institusi,
		Phone:             user.Phone,
		Jenjang:           user.Jenjang,
		Major:             user.Major,
		ProfilePicture:    user.ProfilePicture,
		ProfileHasUpdated: user.ProfileHasUpdated,
		DataHasVerified:   user.DataHasVerified,
		CreatedAt:         user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Handle nullable fields
	if user.NIM != nil {
		userResponse.NIM = *user.NIM
	}
	if user.Gender != nil {
		userResponse.Gender = *user.Gender
	}
	if user.BirthPlace != nil {
		userResponse.BirthPlace = *user.BirthPlace
	}
	if user.BirthDate != nil {
		userResponse.BirthDate = user.BirthDate.Format("2006-01-02")
	}
	if user.IDTeam != nil {
		userResponse.IDTeam = *user.IDTeam
	}
	userResponse.SocMedDocument = user.SocMedDocument
	userResponse.DokumenFilename = user.DokumenFilename

	return userResponse, nil
}

// AdminUpdateUser - Update user by admin
func (s *UserService) AdminUpdateUser(id uint64, updateData dto.AdminUpdateUserDTO) (*dto.AdminUserResponseDTO, error) {
	var user entity.User
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	// Update fields
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Role != "" {
		user.Role = updateData.Role
	}
	if updateData.Institusi != "" {
		user.Institusi = updateData.Institusi
	}
	if updateData.Phone != "" {
		user.Phone = updateData.Phone
	}
	if updateData.Jenjang != "" {
		user.Jenjang = updateData.Jenjang
	}
	if updateData.Major != "" {
		user.Major = updateData.Major
	}
	if updateData.NIM != "" {
		user.NIM = &updateData.NIM
	}
	if updateData.Gender != "" {
		user.Gender = &updateData.Gender
	}
	if updateData.BirthPlace != "" {
		user.BirthPlace = &updateData.BirthPlace
	}
	if updateData.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", updateData.BirthDate)
		if err == nil {
			user.BirthDate = &birthDate
		}
	}
	if updateData.SocMedDocument != "" {
		user.SocMedDocument = updateData.SocMedDocument
	}
	if updateData.DokumenFilename != "" {
		user.DokumenFilename = updateData.DokumenFilename
	}
	if updateData.ProfilePicture != "" {
		user.ProfilePicture = updateData.ProfilePicture
	}

	// Update boolean fields
	user.ProfileHasUpdated = updateData.ProfileHasUpdated
	user.DataHasVerified = updateData.DataHasVerified

	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return s.AdminGetUserById(id)
}

// AdminDeleteUser - Delete user by admin (soft delete)
func (s *UserService) AdminDeleteUser(id uint64) error {
	return s.DB.Delete(&entity.User{}, id).Error
}

// AdminGetUserGrowthAnalytics - Get user growth analytics
func (s *UserService) AdminGetUserGrowthAnalytics(query dto.UserGrowthAnalyticsDTO) ([]dto.UserGrowthResponseDTO, error) {
	startDate, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		return nil, err
	}

	var results []dto.UserGrowthResponseDTO

	// Generate daily analytics between start and end date
	currentDate := startDate
	for currentDate.Before(endDate.Add(24 * time.Hour)) {
		nextDate := currentDate.Add(24 * time.Hour)

		// Count new users for this day
		var newUsers int64
		if err := s.DB.Model(&entity.User{}).Where("created_at >= ? AND created_at < ?", currentDate, nextDate).Count(&newUsers).Error; err != nil {
			return nil, err
		}

		// Count total users up to this day
		var totalUsers int64
		if err := s.DB.Model(&entity.User{}).Where("created_at < ?", nextDate).Count(&totalUsers).Error; err != nil {
			return nil, err
		}

		results = append(results, dto.UserGrowthResponseDTO{
			Period:     currentDate.Format("2006-01-02"),
			NewUsers:   newUsers,
			TotalUsers: totalUsers,
		})

		currentDate = nextDate
	}

	return results, nil
}
