package dto

import (
	"mime/multipart"
	"time"
)

type UserUpdateRequest struct {
	Id        string                `form:"id"`
	FullName  *string               `form:"fullname"`
	Phone     *string               `form:"phone"`
	Email     string                `form:"email" binding:"omitempty,email"`
	AvatarUrl *multipart.FileHeader `form:"avatar_url"`
}

type UserUpdateResponse struct {
	Id        string  `json:"id"`
	FullName  *string `json:"fullname"`
	Phone     *string `json:"phone"`
	Email     string  `json:"email"`
	AvatarUrl *string `json:"avatar_url"`
}

type UserResponse struct {
	ID        string     `json:"id"`
	Fullname  *string    `json:"fullName"`
	Phone     *string    `json:"phone"`
	Avatar    *string    `json:"avatar"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdateAt  *time.Time `json:"updated_at"`
}
type UserChangePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserChangePin struct {
	Pin string `json:"pin"`
}

type SetPinRequest struct {
	Pin string `json:"pin" binding:"required,numeric,len=6"`
}
