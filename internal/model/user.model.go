package model

import "time"

type User struct {
	Id        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Auth struct {
	Id       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	HasPin   bool   `db:"has_pin"`
}

type UserResponse struct {
	ID        string     `db:"id"`
	Fullname  *string    `db:"fullname"`
	Phone     *string    `db:"phone"`
	Avatar    *string    `db:"avatar"`
	Email     string     `db:"email"`
	Pin       *string    `db:"pin"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
