package dto

type CreateNewsLetterDTO struct {
	Title      string `json:"title" binding:"required"`
	NewsLetter string `json:"news_letter" binding:"required"`
	BaseImage  string `json:"base_image"`
	IdAdmin    uint64 `json:"id_admin" binding:"required"`
}

type UpdateNewsLetterDTO struct {
	Title      string `json:"title"`
	NewsLetter string `json:"news_letter"`
	BaseImage  string `json:"base_image"`
}

type NewsLetterResponseDTO struct {
	ID         uint64 `json:"id_news_letter"`
	Title      string `json:"title"`
	NewsLetter string `json:"news_letter"`
	BaseImage  string `json:"base_image"`
	IdAdmin    uint64 `json:"id_admin"`
}
