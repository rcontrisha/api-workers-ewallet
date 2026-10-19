package dto

import "time"

type SearchUserResponse struct {
	ID       string  `json:"id"`
	FullName *string `json:"fullname"`
	Phone    *string `json:"phone"`
	Avatar   *string `json:"avatar"`
}

type TransferRequest struct {
	ReceiverID string  `json:"receiver_id" binding:"required,uuid"`
	Amount     int   `json:"amount" binding:"required,gt=0"`
	Pin        string  `json:"pin" binding:"required,numeric,len=6"`
	Note       *string `json:"note"`
}

type TransferResponse struct {
	TrxCode      string    `json:"trx_code"`
	Amount       int     `json:"amount"`
	Note         *string   `json:"note"`
	Status       string    `json:"status"`
	ReceiverName string    `json:"receiver_name"`
	CreatedAt    time.Time `json:"created_at"`
}