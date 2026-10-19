package model

import "time"

type Transaction struct {
	ID                string     `db:"id"`
	TrxCode           string     `db:"trx_code"`
	Type              string     `db:"type"`
	Status            string     `db:"status"`
	Amount            int      `db:"amount"`
	Tax               int      `db:"tax"`
	Discount          int      `db:"discount"`
	Total             int      `db:"total"`
	SenderWalletID    *string    `db:"sender_wallet_id"`
	ReceiverWalletID  string     `db:"receiver_wallet_id"`
	PaymentMethodID   *string    `db:"payment_method_id"`
	Note              *string    `db:"note"`
	DeletedBySender   bool       `db:"deleted_by_sender"`
	DeletedByReceiver bool       `db:"deleted_by_receiver"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at"`
}