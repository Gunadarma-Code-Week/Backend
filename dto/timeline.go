package dto

type TimelineCreateDTO struct {
	Category    string   `json:"category" binding:"required,max=50"`
	OrderIndex  int      `json:"order_index"`
	Date        string   `json:"date" binding:"required,max=100"`
	Title       string   `json:"title" binding:"required,max=255"`
	Description string   `json:"description"`
	Events      []string `json:"events"`
}

type TimelineUpdateDTO struct {
	Category    string   `json:"category" binding:"omitempty,max=50"`
	OrderIndex  int      `json:"order_index"`
	Date        string   `json:"date" binding:"omitempty,max=100"`
	Title       string   `json:"title" binding:"omitempty,max=255"`
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
