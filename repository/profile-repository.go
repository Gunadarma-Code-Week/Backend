package repository

import (
	"gcw/entity"

	"gorm.io/gorm"
)

type profileRepository struct {
	DB *gorm.DB
}

type ProfileRepository interface {
	Create(*entity.Profile) error
	FindById(uint64, *entity.Profile) error
}

func GateProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{
		DB: db,
	}
}

func (r *profileRepository) Create(u *entity.Profile) error {
	res := r.DB.Create(&u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (r *profileRepository) FindById(id uint64, u *entity.Profile) error {
	res := r.DB.First(&u, id)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}
