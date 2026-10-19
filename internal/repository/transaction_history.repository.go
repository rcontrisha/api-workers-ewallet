package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rcontrisha/api-workers-ewallet/internal/model"
)

type TransactionHistoryRepo struct {
}

func NewTransactionHistoryRepo() *TransactionHistoryRepo {
	return &TransactionHistoryRepo{}
}

func (r *TransactionHistoryRepo) GetHistoryTransaction(ctx context.Context, db DBTX, WalletID, search string, limit, offset int) ([]model.TransactionHistory, error) {
	query := `SELECT
    t.id,
  CASE WHEN t.type = 'topup' THEN pm.name ELSE a.fullname END AS name,
  CASE WHEN t.type = 'topup' THEN pm.logo ELSE a.avatar   END AS avatar,
  CASE WHEN t.type = 'topup' THEN NULL    ELSE a.phone    END AS phone,
  CASE
    WHEN t.type = 'topup'        THEN 'topup'
    WHEN t.sender_wallet_id = $1 THEN 'transfer_out'
    ELSE                              'transfer_in'
  END AS type,
  CASE
    WHEN t.type <> 'topup' AND t.sender_wallet_id = $1 THEN -t.total
    ELSE t.amount
  END AS amount
FROM transaction t
LEFT JOIN account a
  ON a.id = CASE
              WHEN t.sender_wallet_id = $1 THEN t.receiver_wallet_id
              ELSE t.sender_wallet_id
            END
LEFT JOIN payment_method pm ON pm.id = t.payment_method_id
WHERE (
        (t.sender_wallet_id   = $1 AND t.deleted_by_sender   = false)
     OR (t.receiver_wallet_id = $1 AND t.deleted_by_receiver = false)
      )
      AND t.status = 'success'
  AND (
        $2 = ''
     OR CASE WHEN t.type = 'topup' THEN pm.name ELSE a.fullname END ILIKE '%' || $2 || '%'
     OR a.phone ILIKE '%' || $2 || '%'
      )
ORDER BY t.created_at DESC
LIMIT $3 OFFSET $4;`

	args := []any{WalletID, search, limit, offset}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	trxHistories := make([]model.TransactionHistory, 0)
	for rows.Next() {
		var history model.TransactionHistory
		if err := rows.Scan(&history.ID, &history.Name, &history.Avatar, &history.Phone, &history.Type, &history.Amount); err != nil {
			return nil, err
		}
		trxHistories = append(trxHistories, history)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trxHistories, nil
}

func (r *TransactionHistoryRepo) GetTotalDataTransactionHistory(ctx context.Context, db DBTX, walletID, search string) (int, error) {
	query := `
SELECT COUNT(*)
FROM transaction t
LEFT JOIN account a
  ON a.id = CASE
              WHEN t.sender_wallet_id = $1 THEN t.receiver_wallet_id
              ELSE t.sender_wallet_id
            END
LEFT JOIN payment_method pm ON pm.id = t.payment_method_id
WHERE (
        (t.sender_wallet_id   = $1 AND t.deleted_by_sender   = false)
     OR (t.receiver_wallet_id = $1 AND t.deleted_by_receiver = false)
      )
  AND t.status = 'success'
  AND (
        $2::text = ''
     OR CASE WHEN t.type = 'topup' THEN pm.name ELSE a.fullname END ILIKE '%' || $2 || '%'
     OR a.phone ILIKE '%' || $2 || '%'
      )`

	var totalData int
	if err := db.QueryRow(ctx, query, walletID, search).Scan(&totalData); err != nil {
		return 0, err
	}
	return totalData, nil
}

func (r *TransactionHistoryRepo) DeleteHistoryTransaction(ctx context.Context, db DBTX, transactionID, walletID string) error {
	query := `UPDATE transaction
SET deleted_by_sender = CASE
        WHEN sender_wallet_id = $2 THEN TRUE
        ELSE deleted_by_sender
    END,
    deleted_by_receiver = CASE
        WHEN receiver_wallet_id = $2 THEN TRUE
        ELSE deleted_by_receiver
    END
WHERE id = $1
  AND (sender_wallet_id = $2 OR receiver_wallet_id = $2)
  AND (
        (sender_wallet_id = $2 AND deleted_by_sender = FALSE)
     OR (receiver_wallet_id = $2 AND deleted_by_receiver = FALSE)
  )`

	result, err := db.Exec(ctx, query, transactionID, walletID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
