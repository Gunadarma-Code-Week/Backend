package dto

import "time"

type ProfileResponseDTO struct {
	Name       string    `json:"Name"`
	Gender     string    `json:"Gender"`
	NIM        int64     `json:"NIM"`
	Age        int64     `json:"Age"`
	BirthPlace string    `json:"BirthPlace"`
	BirthDate  time.Time `json:"BirthDate"`
	Institusi  string    `json:"Institusi"`
	UserID     string    `json:"UserID"`
}
