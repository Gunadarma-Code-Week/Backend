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

type LoginDTO struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
}

type RegisterDTO struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
	Name     string `json:"name" binding:"required" validate:"min=3"`
}
