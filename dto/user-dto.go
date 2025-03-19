package dto

type UserResponseDTO struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AuthResponseDTO struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	User         UserResponseDTO `json:"user"`
}

// type UserRequestDTO struct {
// 	Email  string `json:"email"`
// 	RoleID uint64 `json:"role_id"`
// 	IDTeam *int   `json:"id_team"`
// }

type ValidateGoogleIdTokenDTO struct {
	GoogleIdToken string `json:"google_id_token"`
}

// type RegisterDTO struct {
// 	GoogleIdToken string `json:"google_id_token"`
// }
