package dto

type TransactionHistoryQuery struct {
	Search string `form:"search"`
	Limit  int    `form:"limit"`
	Page   int    `form:"page"`
}

type TransactionHistory struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
	Phone  *string `json:"phone"`
	Type   string  `json:"type"`
	Amount int64   `json:"amount"`
}
