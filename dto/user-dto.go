package dto

type UserResponseDTO struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

type UserRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
