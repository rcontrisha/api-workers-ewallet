package dto

type DashboardSummaryResponse struct {
	Balance int64 `json:"balance"`
	Income  int64 `json:"income"`
	Expense int64 `json:"expense"`
}

type HistoryTransactionResponse struct {
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
	Type   string  `json:"type"`
	Amount string  `json:"amount"`
}
