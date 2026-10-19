package model

type TransactionHistory struct {
	ID     string  `db:"id"`
	Name   string  `db:"name"`
	Avatar *string `db:"avatar"`
	Phone  *string `db:"phone"`
	Type   string  `db:"type"`
	Amount int64   `db:"amount"`
}
