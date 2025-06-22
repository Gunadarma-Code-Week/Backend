package repository

import (
	"gcw/entity"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

type UserRepository interface {
	FindByUsername(string, *entity.User) error
	FindByEmail(string, *entity.User) error
	FindById(uint64, *entity.User) error
	FindAll(time.Time, time.Time, int, int) ([]*entity.User, int64, error)
	FindByIdTeam(id uint64, users *[]entity.User) error
	Create(*entity.User) error
	Update(u *entity.User, id uint64) error
	UpdateTeamId(idUser uint64, idTeam uint64) error
	GetDB() *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) GetDB() *gorm.DB {
	return r.DB
}

func (r *userRepository) FindByEmail(email string, u *entity.User) error {
	res := r.DB.Where("email = ?", email).First(&u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByUsername(username string, u *entity.User) error {
	res := r.DB.Where("username = ?", username).First(&u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

// find by id
func (r *userRepository) FindById(id uint64, u *entity.User) error {
	res := r.DB.Where("id = ?", id).First(&u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindAll(startDate, endDate time.Time, limit, offset int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var totalUsers int64

	// Query to count the total number of users within the date range
	if err := r.DB.Model(&entity.User{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&totalUsers).Error; err != nil {
		return nil, 0, err
	}

	// Query to fetch the users with pagination and date range filtering
	if err := r.DB.Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalUsers, nil
}

func (r *userRepository) FindByIdTeam(id uint64, users *[]entity.User) error {
	res := r.DB.Where("id_team = ?", id).Find(users)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Create(u *entity.User) error {
	res := r.DB.Create(&u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Update(u *entity.User, id uint64) error {
	u.ID = id
	res := r.DB.Where("id = ?", id).Updates(&u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) UpdateTeamId(idUser uint64, idTeam uint64) error {
	res := r.DB.Model(&entity.User{}).Where("id = ?", idUser).Update("id_team", idTeam)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}
