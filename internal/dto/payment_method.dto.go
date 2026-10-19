package dto

type PaymentMethod struct {
	ID                    string  `json:"id"`
	GatewayPaymentChannel string  `json:"gateway_payment_channel"`
	Name                  string  `json:"name"`
	Logo                  *string `json:"logo"`
	IsActive              bool    `json:"is_active"`
}
