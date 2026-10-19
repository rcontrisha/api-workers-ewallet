package model

type DashboardSummary struct {
	Balance int64 `db:"balance"`
	Income  int64 `db:"income"`
	Expense int64 `db:"expense"`
}

type DashboardHistoryTransaction struct {
	Name   string  `db:"name"`
	Avatar *string `db:"avatar"`
	Type   string  `db:"type"`
	Amount string  `db:"amount"`
}
