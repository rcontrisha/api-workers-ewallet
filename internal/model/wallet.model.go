package model

import "time"

type Wallet struct {
	Id        string    `db:"id"`
	Balance   int       `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
