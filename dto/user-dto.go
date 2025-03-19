package dto

import (
	"time"
)

type UserResponseDTO struct {
	ID         uint64     `json:"id"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	Name       string     `json:"name"`
	Gender     string     `json:"gender"`
	NIM        string     `json:"nim"`
	BirthPlace string     `json:"birth_place"`
	BirthDate  *time.Time `json:"birth_date"`
	Institusi  string     `json:"institusi"`
}

type UpdateUserProfileDTO struct {
	Name       string `json:"name" binding:"required"`
	Gender     string `json:"gender" binding:"required"`
	NIM        string `json:"nim" binding:"required"`
	BirthPlace string `json:"birth_place" binding:"required"`
	BirthDate  string `json:"birth_date" binding:"required"`
	Institusi  string `json:"institusi" binding:"required"`
}

// type RegisterDTO struct {
// 	GoogleIdToken string `json:"google_id_token"`
// }
