package repository

import (
	"context"
	"errors"

	"github.com/rcontrisha/api-workers-ewallet/internal/model"
)

type PaymentRepo struct{}

func NewPaymentRepo() *PaymentRepo {
	return &PaymentRepo{}
}

func (p *PaymentRepo) GetPaymentMethodByID(ctx context.Context, db DBTX, id string) (model.PaymentMethod, error) {
	query := `SELECT id, gateway_payment_channel, name, logo, is_active FROM payment_method WHERE id = $1`
	var pm model.PaymentMethod
	err := db.QueryRow(ctx, query, id).Scan(&pm.ID, &pm.GatewayPaymentChannel, &pm.Name, &pm.Logo, &pm.IsActive)
	if err != nil {
		return model.PaymentMethod{}, err
	}
	return pm, nil
}

func (p *PaymentRepo) CreateTransaction(ctx context.Context, db DBTX, t model.Transaction) (model.Transaction, error) {
	query := `
		INSERT INTO transaction (
			trx_code, type, status, amount, tax, discount, total, 
			sender_wallet_id, receiver_wallet_id, payment_method_id, note
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at`

	err := db.QueryRow(ctx, query,
		t.TrxCode, t.Type, t.Status, t.Amount, t.Tax, t.Discount, t.Total,
		t.SenderWalletID, t.ReceiverWalletID, t.PaymentMethodID, t.Note,
	).Scan(&t.ID, &t.CreatedAt)

	if err != nil {
		return model.Transaction{}, err
	}
	return t, nil
}

func (p *PaymentRepo) FindTransactionByTrxCode(ctx context.Context, db DBTX, trxCode string) (model.Transaction, error) {
	query := `
		SELECT id, trx_code, type, status, amount, tax, discount, total, 
		       sender_wallet_id, receiver_wallet_id, payment_method_id, note, created_at
		FROM transaction WHERE trx_code = $1`

	var t model.Transaction
	err := db.QueryRow(ctx, query, trxCode).Scan(
		&t.ID, &t.TrxCode, &t.Type, &t.Status, &t.Amount, &t.Tax, &t.Discount, &t.Total,
		&t.SenderWalletID, &t.ReceiverWalletID, &t.PaymentMethodID, &t.Note, &t.CreatedAt,
	)
	if err != nil {
		return model.Transaction{}, err
	}
	return t, nil
}

func (p *PaymentRepo) UpdateTransactionStatus(ctx context.Context, db DBTX, trxCode string, status string) error {
	query := `UPDATE transaction SET status = $1, updated_at = NOW() WHERE trx_code = $2`
	cmd, err := db.Exec(ctx, query, status, trxCode)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("transaction not found")
	}
	return nil
}

func (p *PaymentRepo) AddWalletBalance(ctx context.Context, db DBTX, walletID string, amount int) error {
	query := `UPDATE wallet SET balance = balance + $1, updated_at = NOW() WHERE id = $2`
	cmd, err := db.Exec(ctx, query, amount, walletID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("wallet not found")
	}
	return nil
}

func (p *PaymentRepo) GetPaymentMethodList(ctx context.Context, db DBTX) ([]model.PaymentMethod, error) {
	query := `SELECT id, gateway_payment_channel, name, logo, is_active FROM payment_method`
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var listPayments []model.PaymentMethod
	for rows.Next() {
		var payment model.PaymentMethod
		if err := rows.Scan(&payment.ID, &payment.GatewayPaymentChannel, &payment.Name, &payment.Logo, &payment.IsActive); err != nil {
			return nil, err
		}
		listPayments = append(listPayments, payment)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return listPayments, nil
}
