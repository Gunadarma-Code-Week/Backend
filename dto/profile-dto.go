package dto

type ProfileResponseDTO struct {
	Name       string    `json:"Name"`
	NIM        int64     `json:"NIM"`
	Age        int64     `json:"Age"`
	Institusi  string    `json:"Institusi"`
	UserID     string    `json:"UserID"`
}
