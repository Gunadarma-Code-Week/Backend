package dto

import "time"

type User struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	University string `json:"university"`
	Degree     string `json:"degree"`
}

type Member struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

type Team struct {
	TeamName string   `json:"team_name"`
	Members  []Member `json:"members"`
}

type Ticket struct {
	TicketId   string    `json:"ticket_id"`
	Type       string    `json:"type"`
	IssuedAt   time.Time `json:"issued_at"`
	ValidUntil time.Time `json:"valid_until"`
	QrCodeUrl  string    `json:"qr_code_url"`
}

type Event struct {
	Name          string `json:"name"`
	Status        string `json:"status"`
	PaymentStatus string `json:"payment_status"`
	Team          []Team `json:"team"`
	IdTeam        int64  `json:"id_team"`
	Ticket        Ticket `json:"ticket"`
}

type ResponseEvents struct {
	User   User    `json:"user"`
	Events []Event `json:"events"`
}
