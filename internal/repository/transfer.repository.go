package repository

import (
	"context"
	"errors"
	"log"

	"github.com/rcontrisha/api-workers-ewallet/internal/model"
)

type TransferRepo struct{}

func NewTransferRepo() *TransferRepo {
	return &TransferRepo{}
}

func (t *TransferRepo) GetRecentAndFavoriteContacts(ctx context.Context, db DBTX, currentUserID string) ([]model.Account, error) {
	query := `
		SELECT DISTINCT ON (u.id) u.id, a.fullname, a.phone, a.avatar
		FROM (
			SELECT receiver_wallet_id AS user_id, created_at
			FROM transaction
			WHERE sender_wallet_id = $1 AND type = 'transfer' AND status = 'success'
			
			UNION ALL
			
			SELECT favorite_user_id AS user_id, NOW() AS created_at
			FROM favorite
			WHERE user_id = $1
		) c
		JOIN users u ON u.id = c.user_id
		JOIN account a ON a.id = u.id
		WHERE u.id != $1
		ORDER BY u.id, c.created_at DESC
		LIMIT 10`

	rows, err := db.Query(ctx, query, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []model.Account
	for rows.Next() {
		var acc model.Account
		if err := rows.Scan(&acc.Id, &acc.FullName, &acc.Phone, &acc.Avatar); err != nil {
			return nil, err
		}
		contacts = append(contacts, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

// SearchUsers melakukan pencarian user global berdasarkan nama, nomor hp, atau email
func (t *TransferRepo) SearchUsers(ctx context.Context, db DBTX, query string, excludeUserID string) ([]model.Account, error) {
	sqlQuery := `
		SELECT u.id, a.fullname, a.phone, a.avatar
		FROM users u
		JOIN account a ON a.id = u.id
		WHERE u.id != $1 AND (a.fullname ILIKE $2 OR a.phone ILIKE $2 OR u.email ILIKE $2)
		LIMIT 20`

	searchPattern := "%" + query + "%"
	rows, err := db.Query(ctx, sqlQuery, excludeUserID, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.Account
	for rows.Next() {
		var acc model.Account
		if err := rows.Scan(&acc.Id, &acc.FullName, &acc.Phone, &acc.Avatar); err != nil {
			return nil, err
		}
		users = append(users, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (t *TransferRepo) GetUserPINAndProfile(ctx context.Context, db DBTX, userID string) (string, model.Account, error) {
	query := `SELECT id, fullname, phone, avatar, pin FROM account WHERE id = $1`
	var acc model.Account
	var pin *string

	err := db.QueryRow(ctx, query, userID).Scan(&acc.Id, &acc.FullName, &acc.Phone, &acc.Avatar, &pin)
	if err != nil {
		log.Println(err)
		return "", model.Account{}, err
	}

	if pin == nil || *pin == "" {
		return "", acc, errors.New("PIN_NOT_SET")
	}

	return *pin, acc, nil
}

func (t *TransferRepo) DeductWalletBalance(ctx context.Context, db DBTX, walletID string, amount int) error {
	query := `UPDATE wallet SET balance = balance - $1, updated_at = NOW() WHERE id = $2 AND balance >= $1`
	cmd, err := db.Exec(ctx, query, amount, walletID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("INSUFFICIENT_BALANCE")
	}
	return nil
}

func (t *TransferRepo) AddWalletBalance(ctx context.Context, db DBTX, walletID string, amount int) error {
	query := `UPDATE wallet SET balance = balance + $1, updated_at = NOW() WHERE id = $2`
	cmd, err := db.Exec(ctx, query, amount, walletID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("RECEIVER_WALLET_NOT_FOUND")
	}
	return nil
}

func (t *TransferRepo) CreateTransaction(ctx context.Context, db DBTX, tr model.Transaction) (model.Transaction, error) {
	query := `
		INSERT INTO transaction (
			trx_code, type, status, amount, tax, discount, total,
			sender_wallet_id, receiver_wallet_id, note
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`

	err := db.QueryRow(ctx, query,
		tr.TrxCode, tr.Type, tr.Status, tr.Amount, tr.Tax, tr.Discount, tr.Total,
		tr.SenderWalletID, tr.ReceiverWalletID, tr.Note,
	).Scan(&tr.ID, &tr.CreatedAt)

	if err != nil {
		return model.Transaction{}, err
	}
	return tr, nil
}
