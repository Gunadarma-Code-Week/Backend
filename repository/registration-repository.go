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
	CreateTeam(*entity.HackathonTeam) error
	UpdateUserTeam(*entity.User, uint64, uint64) error
	UpdateUserTeamJoinCode(*entity.User, string, uint64) error
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

func (r *registrationRepository) CreateTeam(u *entity.HackathonTeam) error {
	res := r.DB.Create(&u)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}

func (r *registrationRepository) UpdateUserTeam(u *entity.User, id_team uint64, id_user uint64) error {
	if err := r.DB.First(u, id_user).Error; err != nil {
		return err
	}

	u.IDTeam = &id_team

	if err := r.DB.Save(u).Error; err != nil {
		return err
	}

	return nil
}

func (r *registrationRepository) UpdateUserTeamJoinCode(u *entity.User, code string, id_user uint64) error {
	team := &entity.Team{}

	if err := r.DB.Where("join_code = ?", code).First(team).Error; err != nil {
		return err
	}

	if err := r.DB.Model(u).
		Where("id = ?", id_user).
		Update("id_team", team.ID_Team).Error; err != nil {
		return err
	}

	return nil
}
