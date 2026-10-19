package model

type PaymentMethod struct {
	ID                    string  `db:"id"`
	GatewayPaymentChannel string  `db:"gateway_payment_channel"`
	Name                  string  `db:"name"`
	Logo                  *string `db:"logo"`
	IsActive              bool    `db:"is_active"`
}
