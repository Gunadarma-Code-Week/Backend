package dto

type AuthResponseDTO struct {
	User UserResponseDTO `json:"user"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ValidateGoogleIdTokenDTO struct {
	GoogleIdToken string `json:"google_id_token" binding:"required"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
