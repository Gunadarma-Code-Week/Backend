package dto

import (
	"time"
)

type UserResponseDTO struct {
	ID             uint64     `json:"id"`
	Email          string     `json:"email"`
	Role           string     `json:"role"`
	Name           string     `json:"name"`
	Major          string     `json:"major"`
	ProfilePicture string     `json:"profile_picture"`
	NIM            string     `json:"nim"`
	BirthDate      *time.Time `json:"birth_date"`
	Institusi      string     `json:"institusi"`
	Phone          string     `json:"phone"`

	SocMedDocument  string `json:"socmed_document"`
	DokumenFilename string `json:"dokumen_filename"`

	ProfileHasUpdated bool `json:"profile_has_updated"`
	DataHasVerified   bool `json:"data_has_verified"`
}

type UpdateUserProfileDTO struct {
	Name       string `json:"name" binding:"required"`
	Gender     string `json:"gender" binding:"required"`
	NIM        string `json:"nim" binding:"required"`
	Phone      string `json:"phone"`
	Major      string `json:"major"`
	BirthPlace string `json:"birth_place" binding:"required"`
	BirthDate  string `json:"birth_date" binding:"required"`
	Institusi  string `json:"institusi" binding:"required"`

	SocMedDocument string `json:"socmed_document"`
}

// type RegisterDTO struct {
// 	GoogleIdToken string `json:"google_id_token"`
// }
