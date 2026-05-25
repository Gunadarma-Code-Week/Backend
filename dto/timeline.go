package dto

type TimelineCreateDTO struct {
	Category    string   `json:"category" binding:"required"`
	OrderIndex  int      `json:"order_index"`
	Date        string   `json:"date" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Events      []string `json:"events"`
}

type TimelineUpdateDTO struct {
	Category    string   `json:"category"`
	OrderIndex  int      `json:"order_index"`
	Date        string   `json:"date"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Events      []string `json:"events"`
}

type TimelineResponseDTO struct {
	ID          uint     `json:"id"`
	Category    string   `json:"category"`
	OrderIndex  int      `json:"order_index"`
	Date        string   `json:"date"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Events      []string `json:"events"`
}
