package dto

import "time"

type RegistraionTeamRequest struct {
	TeamName       string `json:"team_name" binding:"required"`
	KomitmenFee    string `json:"bukti_pembayran"`
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
	KomitmenFee      string `json:"bukti_pembayaran"`
	JoinCode         string `json:"join_code"`
	RegistrationHackathonRequest
	IDTeam    uint64    `json:"id_team"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegistrationCPResponse struct {
	ID_CPTeam        uint64 `json:"id_cp_team"`
	Stage            string `json:"stage"`
	Status           string `json:"status"`
	KomitmenFee      string `json:"bukti_pembayaran"`
	DomjudgeUsername string `json:"domjudge_username"`
	DomjudgePassword string `json:"domjudge_password"`
	JoinCode         string `json:"join_code"`
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
