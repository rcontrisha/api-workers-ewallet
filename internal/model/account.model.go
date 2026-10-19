package model

import "time"

type Account struct {
	Id         string    `db:"id"`
	Email      string    `db:"email"`
	FullName   *string    `db:"fullname"`
	Phone      *string    `db:"phone"`
	Avatar     *string   `db:"avatar"`
	Pin        string    `db:"pin"`
	IsVerified bool      `db:"is_verfied"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
