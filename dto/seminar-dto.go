package dto

import "time"

// JoinSeminarRequest DTO untuk request join seminar
type JoinSeminarRequest struct {
	// ID tiket akan di-generate otomatis oleh sistem
}

// JoinSeminarResponse DTO untuk response join seminar
type JoinSeminarResponse struct {
	Message  string `json:"message" binding:"max=50" example:"Berhasil bergabung ke seminar"`
	Status   string `json:"status" binding:"max=50" example:"success"`
	IDTiket  string `json:"id_tiket" binding:"max=50" example:"TICKET123456"`
	SeminarID uint64 `json:"seminar_id" example:"1"`
}

// Response DTO untuk detail tiket seminar
type SeminarTicketDetail struct {
	IDSeminar     uint64    `json:"id_seminar"`
	IDTiket       string    `json:"id_tiket" binding:"max=50"`
	PaymentStatus string    `json:"payment_status" binding:"max=50"`
	User          UserInfo  `json:"user"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserInfo struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name" binding:"max=50"`
	Email    string `json:"email" binding:"max=100"`
	Phone    string `json:"phone" binding:"max=50"`
	Jenjang  string `json:"jenjang" binding:"max=50"`
	Institusi string `json:"institusi" binding:"max=50"`
}

// AdminAddParticipantRequest DTO untuk admin menambahkan participant ke seminar
type AdminAddParticipantRequest struct {
	UserID uint64 `json:"user_id" binding:"required" example:"1"`
}

// AdminAddParticipantResponse DTO untuk response admin menambahkan participant
type AdminAddParticipantResponse struct {
	Message   string `json:"message" binding:"max=50" example:"Berhasil menambahkan participant ke seminar"`
	Status    string `json:"status" binding:"max=50" example:"success"`
	IDTiket   string `json:"id_tiket" binding:"max=50" example:"SEM20241201ABC123"`
	SeminarID uint64 `json:"seminar_id" example:"1"`
	User      UserInfo `json:"user"`
}