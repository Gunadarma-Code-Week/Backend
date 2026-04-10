package dto

import "time"

type RegistraionTeamRequest struct {
	TeamName       string `json:"team_name" binding:"required"`
	Supervisor     string `json:"supervisor" binding:"required"`
	SupervisorNIDN string `json:"supervisor_nidn" binding:"required"`
}

type RegistrationHackathonRequest struct {
	BuktiPembayaran string `json:"bukti_pembayaran"`
}

type RegistrationHackathonTeamRequest struct {
	RegistraionTeamRequest
	RegistrationHackathonRequest
}

type RegistrationCPRequest struct {
	JoinCode        string `json:"join_code"`
	BuktiPembayaran string `json:"bukti_pembayaran"`
}

type RegistrationCPTeamRequest struct {
	RegistraionTeamRequest
	RegistrationCPRequest
}

type RegistraionTeamResponse struct {
	ID_Team uint64   `json:"id_team"`
	Members []Member `json:"member"`
	Leader  Member   `json:"leader"`
	QRString string   `json:"qr_string,omitempty"`
	OrderID  string   `json:"order_id,omitempty"`
	PaymentStatus string `json:"payment_status,omitempty"`
	RegistraionTeamRequest
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegistrationHackathonResponse struct {
	ID_HackathonTeam uint64 `gorm:"primary_key:auto_increment"`
	Stage            string `gorm:"varchar(255); not null"`
	Status           string `gorm:"varchar(255); not null"`
	JoinCode         string `json:"join_code"`
	QRString         string `json:"qr_string,omitempty"`
	OrderID          string `json:"order_id,omitempty"`
	PaymentStatus    string `json:"payment_status,omitempty"`
	KomitmenFee      string `json:"bukti_pembayaran,omitempty"`
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
	JoinCode         string `json:"join_code"`
	QRString         string `json:"qr_string,omitempty"`
	OrderID          string `json:"order_id,omitempty"`
	PaymentStatus    string `json:"payment_status,omitempty"`
	KomitmenFee      string `json:"bukti_pembayaran,omitempty"`
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

type RegistrationCTFRequest struct {
}

type RegistrationCTFTeamRequest struct {
	RegistraionTeamRequest
	RegistrationCTFRequest
	BuktiPembayaran string `json:"bukti_pembayaran"`
}

type RegistrationCTFResponse struct {
	ID_CTFTeam    uint64    `json:"id_ctf_team"`
	Stage         string    `json:"stage"`
	Status        string    `json:"status"`
	JoinCode      string    `json:"join_code"`
	QRString      string    `json:"qr_string,omitempty"`
	OrderID       string    `json:"order_id,omitempty"`
	PaymentStatus string    `json:"payment_status,omitempty"`
	KomitmenFee   string    `json:"bukti_pembayaran,omitempty"`
	IDTeam        uint64    `json:"id_team"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type RegistrationCTFTeamResponse struct {
	Team    RegistraionTeamResponse
	CTFTeam RegistrationCTFResponse
}
