package dto

type ProfileResponseDTO struct {
	Name       string    `json:"Name" binding:"max=50"`
	NIM        int64     `json:"NIM"`
	Age        int64     `json:"Age"`
	Institusi  string    `json:"Institusi" binding:"max=50"`
	UserID     string    `json:"UserID" binding:"max=50"`
}
