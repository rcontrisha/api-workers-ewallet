package repository

import (
	"context"
	"errors"
	"log"

	msgerror "github.com/rcontrisha/api-workers-ewallet/internal/error"
	"github.com/rcontrisha/api-workers-ewallet/internal/model"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (u *UserRepo) UpdateProfileInfo(ctx context.Context, userId string, db DBTX, payload model.Account) (model.Account, error) {
	sql := `UPDATE account
					SET fullname = COALESCE(NULLIF($1, ''),fullname),
							phone = COALESCE(NULLIF($2,''),phone),
							avatar = COALESCE(NULLIF($3, ''),avatar),
							updated_at = now()
					WHERE id = $4
					RETURNING
					fullname, phone, avatar, id`

	var account model.Account

	if err := db.QueryRow(ctx, sql,
		payload.FullName,
		payload.Phone,
		payload.Avatar,
		userId).Scan(&account.FullName, &account.Phone, &account.Avatar, &account.Id); err != nil {
		return model.Account{}, err
	}

	return account, nil

}

func (u *UserRepo) UpdateUserEmail(ctx context.Context, email, userId string, db DBTX) (string, error) {
	sql := `UPDATE users
				SET email = COALESCE(NULLIF($1,''),email)
				WHERE id = $2
				RETURNING email`

	var userEmail string
	if err := db.QueryRow(ctx, sql, email, userId).Scan(&userEmail); err != nil {
		return "", err
	}
	return userEmail, nil
}

func (u *UserRepo) GetProfile(ctx context.Context, db DBTX, userID string) (model.UserResponse, error) {
	query := `
		SELECT u.id, a.fullname, a.phone, a.avatar, u.email, a.pin, a.created_at, a.updated_at
		FROM users u
		JOIN account a ON a.id = u.id
		WHERE u.id = $1`
	var user model.UserResponse
	err := db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Fullname,
		&user.Phone,
		&user.Avatar,
		&user.Email,
		&user.Pin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.UserResponse{}, err
	}
	return user, nil
}

func (u *UserRepo) DeleteAvatar(ctx context.Context, userID string, db DBTX) error {
	query := `UPDATE account SET avatar = NULL, updated_at = now() WHERE id = $1`
	result, err := db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return msgerror.ErrNoRows
	}
	return nil
}

func (u *UserRepo) ChangePassword(ctx context.Context, userId string, password string, db DBTX) error {
	sql := `
        UPDATE users
        SET password = $1
        WHERE id = $2
    `

	tag, err := db.Exec(ctx, sql, password, userId)
	if err != nil {
		log.Println("[Repo.ChangePassword] err:", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return msgerror.ErrNoRows
	}

	return nil
}

func (u *UserRepo) GetUserPassword(ctx context.Context, userId string, db DBTX) (string, error) {
	sql := `SELECT password FROM users WHERE id = $1`
	var userPass string
	if err := db.QueryRow(ctx, sql, userId).Scan(&userPass); err != nil {
		return "", err
	}
	return userPass, nil

}

func (u *UserRepo) ChangePinUser(ctx context.Context, userId string, db DBTX, pin string) error {
	sql := `UPDATE account
			SET pin = $1
			WHERE id = $2`
	cmdTag, err := db.Exec(ctx, sql, pin, userId)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return msgerror.ErrNoRows
	}
	return nil

}

func (u *UserRepo) SetUserPIN(ctx context.Context, db DBTX, userID string, hashedPIN string) error {
	query := `UPDATE account SET pin = $1 WHERE id = $2`
	cmd, err := db.Exec(ctx, query, hashedPIN, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("USER_NOT_FOUND")
	}
	return nil
}

func (u *UserRepo) UserSummaryDashboard(ctx context.Context, userId string, db DBTX) (model.DashboardSummary, error) {
	query := `SELECT 
				w.balance, 
				COALESCE(SUM(t.total) FILTER(WHERE t.sender_wallet_id = w.id),0) AS expense, 
				COALESCE(SUM(t.amount) FILTER (WHERE t.receiver_wallet_id = w.id),0) AS income
				FROM wallet w
				LEFT JOIN transaction t ON (t.sender_wallet_id = w.id OR t.receiver_wallet_id = w.id) AND t.status = 'success'
				WHERE w.id = $1
				GROUP BY w.id, w.balance`

	var data model.DashboardSummary
	if err := db.QueryRow(ctx, query, userId).Scan(&data.Balance, &data.Expense, &data.Income); err != nil {
		return model.DashboardSummary{}, err
	}
	return data, nil
}

func (u *UserRepo) UserHistoryTransaction(ctx context.Context, userId string, db DBTX) ([]model.DashboardHistoryTransaction, error) {
	query := `SELECT
  CASE WHEN t.type = 'topup' THEN pm.name ELSE a.fullname END AS name,
  CASE WHEN t.type = 'topup' THEN pm.logo ELSE a.avatar   END AS avatar,
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
ORDER BY t.created_at DESC
LIMIT 10`

	row, err := db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var histories []model.DashboardHistoryTransaction
	for row.Next() {
		var data model.DashboardHistoryTransaction
		if err := row.Scan(&data.Name, &data.Avatar, &data.Type, &data.Amount); err != nil {
			return nil, err
		}
		histories = append(histories, data)

	}
	if err := row.Err(); err != nil {
		return nil, err
	}
	return histories, nil
}
