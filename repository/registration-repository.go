package repository

import (
	"gcw/entity"

	"gorm.io/gorm"
)

type RegistrationRepository struct {
	DB *gorm.DB
}

func GateRegistrationRepository(db *gorm.DB) *RegistrationRepository {
	return &RegistrationRepository{
		DB: db,
	}
}

func (r *RegistrationRepository) CreateTeam(tx *gorm.DB, u *entity.Team) error {
	res := tx.Create(&u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (r *RegistrationRepository) CreateHackathonTeam(tx *gorm.DB, u *entity.HackathonTeam) error {
	res := tx.Create(&u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (r *RegistrationRepository) CreateCPTeam(tx *gorm.DB, u *entity.CPTeam) error {
	res := tx.Create(&u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (r *RegistrationRepository) FindTeamByJoinCode(team *entity.Team, joinCode string) error {
	res := r.DB.Where("join_code = ?", joinCode).First(&team)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (r *RegistrationRepository) UpdateUserTeam(tx *gorm.DB, u *entity.User, id_team uint64, id_user uint64) error {
	res := tx.Model(&u).Where("id = ?", id_user).Update("id_team", id_team)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}
