package dto

type AuthResponseDTO struct {
	User UserResponseDTO `json:"user"`
}

type ValidateGoogleIdTokenDTO struct {
	GoogleIdToken string `json:"google_id_token" binding:"required"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
