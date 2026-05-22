package dto

import (
	"time"
)

type User struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name" binding:"max=50"`
	Email      string `json:"email" binding:"max=100"`
	University string `json:"university" binding:"max=50"`
}

type Member struct {
	Name       string `json:"name" binding:"max=50"`
	Role       string `json:"role" binding:"max=50"`
	Email      string `json:"email" binding:"max=100"`
	University string `json:"university" binding:"max=50"`
}

type Team struct {
	TeamName string   `json:"team_name" binding:"max=50"`
	Members  []Member `json:"members"`
}

type Ticket struct {
	TicketId   string    `json:"ticket_id" binding:"max=50"`
	Type       string    `json:"type" binding:"max=50"`
	IssuedAt   time.Time `json:"issued_at"`
	ValidUntil time.Time `json:"valid_until"`
	QrCodeUrl  string    `json:"qr_code_url"`
}

type Event struct {
	Name        string `json:"name" binding:"max=50"`
	Status      string `json:"status" binding:"max=50"`
	PaymentType string `json:"payment_type" binding:"max=50"`
	Ticket      Ticket `json:"ticket"`
}

type ResponseEvents struct {
	User   User    `json:"user"`
	IdTeam string  `json:"id_team" binding:"max=50"`
	Events []Event `json:"events"`
}
