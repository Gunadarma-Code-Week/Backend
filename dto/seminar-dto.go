package dto

import "time"

// JoinSeminarRequest DTO untuk request join seminar
type JoinSeminarRequest struct {
	// ID tiket akan di-generate otomatis oleh sistem
}

// JoinSeminarResponse DTO untuk response join seminar
type JoinSeminarResponse struct {
	Message  string `json:"message" example:"Berhasil bergabung ke seminar"`
	Status   string `json:"status" example:"success"`
	IDTiket  string `json:"id_tiket" example:"TICKET123456"`
	SeminarID uint64 `json:"seminar_id" example:"1"`
}

// Response DTO untuk detail tiket seminar
type SeminarTicketDetail struct {
	IDSeminar     uint64    `json:"id_seminar"`
	IDTiket       string    `json:"id_tiket"`
	PaymentStatus string    `json:"payment_status"`
	User          UserInfo  `json:"user"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserInfo struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Jenjang  string `json:"jenjang"`
	Institusi string `json:"institusi"`
}