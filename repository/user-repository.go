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
	Create(*entity.User) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) FindByUsername(username string, u *entity.User) error {
	res := r.DB.Where("username = ?", username).First(&u)
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
