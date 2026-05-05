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

func (r *RegistrationRepository) CreateCTFTeam(tx *gorm.DB, u *entity.CTFTeam) error {
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

func (r *RegistrationRepository) CountUserByTeamID(id_team uint64) (int64, error) {
	var count int64
	res := r.DB.Model(&entity.User{}).Where("id_team = ?", id_team).Count(&count)
	if err := res.Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindActiveTeamByNameGlobal checks whether any ACTIVE (non-soft-deleted) team
// across hackathon, cp, or ctf events already uses this name (case-insensitive).
// Returns nil when a conflicting team IS found (caller should reject the name).
// Returns an error (e.g. gorm.ErrRecordNotFound) when the name is available.
func (r *RegistrationRepository) FindActiveTeamByNameGlobal(name string) error {
	var count int64
	err := r.DB.Raw(`
		SELECT COUNT(*) FROM teams
		LEFT JOIN hackathon_teams ON hackathon_teams.id_team = teams.id_team AND teams.event = 'hackathon'
		LEFT JOIN cp_teams        ON cp_teams.id_team        = teams.id_team AND teams.event = 'cp'
		LEFT JOIN ctf_teams       ON ctf_teams.id_team       = teams.id_team AND teams.event = 'ctf'
		WHERE LOWER(teams.team_name) = LOWER(?)
		AND (
			(teams.event = 'hackathon' AND hackathon_teams.is_deleted = false)
			OR (teams.event = 'cp'        AND cp_teams.is_deleted        = false)
			OR (teams.event = 'ctf'       AND ctf_teams.is_deleted       = false)
		)
	`, name).Scan(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // nil = conflict found
	}
	return gorm.ErrRecordNotFound // error = name is available (no conflict)
}


func (r *RegistrationRepository) UpdateUserTeam(tx *gorm.DB, u *entity.User, id_team uint64, id_user uint64) error {
	res := tx.Model(&entity.User{}).Where("id = ?", id_user).Update("id_team", id_team)
	if err := res.Error; err != nil {
		return err
	}
	return nil
}
