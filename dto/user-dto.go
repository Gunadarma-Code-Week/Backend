package dto

type UserResponseDTO struct {
	ID          uint64 `json:"id"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type UserRequestDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   uint64 `json:"role_id"`
	IDTeam   *int   `json:"id_team"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
