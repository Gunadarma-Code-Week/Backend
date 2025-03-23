package dto

import "time"

type RegistraionTeamRequest struct {
	TeamName       string `json:"team_name" binding:"required"`
	Supervisor     string `json:"supervisor" binding:"required"`
	SupervisorNIDN string `json:"supervisor_nidn" binding:"required"`
}

type RegistrationHackathonRequest struct {
}

type RegistrationHackathonTeamRequest struct {
	RegistraionTeamRequest
	RegistrationHackathonRequest
}

type RegistrationCPRequest struct {
}

type RegistrationCPTeamRequest struct {
	RegistraionTeamRequest
	RegistrationCPRequest
}

type RegistraionTeamResponse struct {
	ID_Team uint64 `json:"id_team"`
	RegistraionTeamRequest
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegistrationHackathonResponse struct {
	ID_HackathonTeam uint64 `gorm:"primary_key:auto_increment"`
	Stage            string `gorm:"varchar(255); not null"`
	Status           string `gorm:"varchar(255); not null"`
	RegistrationHackathonRequest
	IDTeam    uint64    `json:"id_team"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegistrationCPResponse struct {
	ID_CPTeam        uint64 `json:"id_cp_team"`
	Stage            string `json:"stage"`
	Status           string `json:"status"`
	DomjudgeUsername string `json:"domjudge_username"`
	DomjudgePassword string `json:"domjudge_password"`
	RegistrationCPRequest
	IDTeam    uint64 `json:"id_team"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RegistrationCPTeamResponse struct {
	Team   RegistraionTeamResponse
	CPTeam RegistrationCPResponse
}

type RegistrationHackathonTeamResponse struct {
	Team          RegistraionTeamResponse
	HackathonTeam RegistrationHackathonResponse
}
