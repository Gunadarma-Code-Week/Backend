package repository

import (
	"gcw/entity"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

type UserRepository interface {
	FindByUsername(string, *entity.User) error
	FindByEmail(string, *entity.User) error
	FindById(uint64, *entity.User) error
	Create(*entity.User) error
	Update(u *entity.User, id uint64) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
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
