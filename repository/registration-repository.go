package repository

import (
	"gcw/entity"

	"gorm.io/gorm"
)

type registrationRepository struct {
	DB *gorm.DB
}

type RegistrationRepository interface {
	Create(*entity.Team) error
}

func GateRegistrationRepository(db *gorm.DB) RegistrationRepository {
	return &registrationRepository{
		DB: db,
	}
}

func (r *registrationRepository) Create(u *entity.Team) error {
	res := r.DB.Create(&u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}
