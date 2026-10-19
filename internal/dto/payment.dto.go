package dto

import "time"

type CreateTopupRequest struct {
	Amount          int     `json:"amount" binding:"required,gt=0"`
	Tax             int     `json:"tax"`
	PaymentMethodID *string `json:"payment_method_id"`
	Note            *string `json:"note"`
	Pin             string  `json:"pin" binding:"required,len=6"`
}

type CreateTopupResponse struct {
	TrxCode    string `json:"trx_code"`
	Amount     int    `json:"amount"`
	Total      int    `json:"total"`
	Status     string `json:"status"`
	PaymentURL string `json:"payment_url"`
}

// Struct penampung Notification Webhook Callback dari DOKU
type DokuNotificationPayload struct {
	Order struct {
		InvoiceNumber string `json:"invoice_number"`
		Amount        int    `json:"amount"`
	} `json:"order"`
	Transaction struct {
		Status string `json:"status"`
	} `json:"transaction"`
}

type TopupSummaryResponse struct {
	TrxCode          string    `json:"trx_code"`
	Status           string    `json:"status"` // PENDING | SUCCESS | FAILED
	Amount           int       `json:"amount"`
	Tax              int       `json:"tax"`
	Total            int       `json:"total"`
	PaymentMethod    *string   `json:"payment_method"`
	PaymentLogo      *string   `json:"payment_logo"`
	DestinationPhone *string   `json:"destination_phone"`
	Note             *string   `json:"note"`
	CreatedAt        time.Time `json:"created_at"`
}
