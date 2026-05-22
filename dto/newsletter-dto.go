package dto

type CreateNewsLetterDTO struct {
	Title      string `json:"title" binding:"required,max=50"`
	NewsLetter string `json:"news_letter" binding:"required,max=50"`
	BaseImage  string `json:"base_image" binding:"max=50"`
	IdAdmin    uint64 `json:"id_admin" binding:"required"`
}

type UpdateNewsLetterDTO struct {
	Title      string `json:"title" binding:"max=50"`
	NewsLetter string `json:"news_letter" binding:"max=50"`
	BaseImage  string `json:"base_image" binding:"max=50"`
}

type NewsLetterResponseDTO struct {
	ID         uint64 `json:"id_news_letter"`
	Title      string `json:"title" binding:"max=50"`
	NewsLetter string `json:"news_letter" binding:"max=50"`
	BaseImage  string `json:"base_image" binding:"max=50"`
	IdAdmin    uint64 `json:"id_admin"`
}
