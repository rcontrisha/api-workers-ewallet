package repository

import (
	"context"
	"log"

	"github.com/rcontrisha/api-workers-ewallet/internal/model"
)

type TransactionHistoryRepo struct {
}

func NewTransactionHistoryRepo() *TransactionHistoryRepo {
	return &TransactionHistoryRepo{}
}

func (r *TransactionHistoryRepo) GetHistoryTransaction(ctx context.Context, db DBTX, WalletID, search string, limit, offset int) ([]model.TransactionHistory, error) {
	log.Println("Limit ", limit)
	log.Println("Search ", search)
	log.Println("Offset ", offset)

	query := `SELECT
  CASE WHEN t.type = 'TOPUP' THEN pm.name ELSE a.fullname END AS name,
  CASE WHEN t.type = 'TOPUP' THEN pm.logo ELSE a.avatar   END AS avatar,
  CASE WHEN t.type = 'TOPUP' THEN NULL    ELSE a.phone    END AS phone,
  CASE
    WHEN t.type = 'TOPUP'        THEN 'TOPUP'
    WHEN t.sender_wallet_id = $1 THEN 'transfer_out'
    ELSE                              'transfer_in'
  END AS type,
  CASE
    WHEN t.type <> 'TOPUP' AND t.sender_wallet_id = $1 THEN -t.total
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
      AND t.status = 'SUCCESS'
  AND (
        $2 = ''
     OR CASE WHEN t.type = 'TOPUP' THEN pm.name ELSE a.fullname END ILIKE '%' || $2 || '%'
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
		if err := rows.Scan(&history.Name, &history.Avatar, &history.Phone, &history.Type, &history.Amount); err != nil {
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
     OR CASE WHEN t.type = 'TOPUP' THEN pm.name ELSE a.fullname END ILIKE '%' || $2 || '%'
     OR a.phone ILIKE '%' || $2 || '%'
      )`

	var totalData int
	if err := db.QueryRow(ctx, query, walletID, search).Scan(&totalData); err != nil {
		return 0, err
	}
	return totalData, nil
}
