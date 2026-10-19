package dto

type Response struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    any                `json:"data"`
	Meta    PaginationMetaData `json:"metadata,omitzero"`
}

type PaginationMetaData struct {
	CurrentPage int     `json:"current_page"`
	TotalPage   float64 `json:"total_page"`
	TotalData   int     `json:"total_data"`
	Next        string  `json:"next"`
	Prev        string  `json:"prev"`
}
