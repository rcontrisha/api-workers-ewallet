package repository

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rcontrisha/api-workers-ewallet/internal/model"
	"github.com/redis/go-redis/v9"
)

type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type AuthRepo struct {
	rdb *redis.Client
}

func NewAuthRepo(rdb *redis.Client) *AuthRepo {
	return &AuthRepo{
		rdb: rdb,
	}
}

func (a *AuthRepo) FindAccount(ctx context.Context, db DBTX, email string) (model.Auth, error) {
	query := `
		SELECT u.id, u.email, u.password, (a.pin IS NOT NULL) AS has_pin
		FROM users u
		LEFT JOIN account a ON a.id = u.id
		WHERE u.email = $1`

	var acc model.Auth
	err := db.QueryRow(ctx, query, email).Scan(&acc.Id, &acc.Email, &acc.Password, &acc.HasPin)
	if err != nil {
		return model.Auth{}, err
	}
	return acc, nil
}

func (a *AuthRepo) CreateUser(ctx context.Context, db DBTX, payload model.User) (string, error) {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	args := []any{payload.Email, payload.Password}

	var newUserId string
	err := db.QueryRow(ctx, query, args...).Scan(
		&newUserId,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}

		return "", err
	}

	return newUserId, nil
}

func (a *AuthRepo) CreateAccount(ctx context.Context, db DBTX, userId string) error {
	query := "INSERT INTO account (id) VALUES ($1)"

	cmd, err := db.Exec(ctx, query, userId)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("No Row(s) Affected")
	}

	return nil
}

func (a *AuthRepo) CreateWallet(ctx context.Context, db DBTX, userId string) error {
	query := "INSERT INTO wallet (id) VALUES ($1)"

	cmd, err := db.Exec(ctx, query, userId)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("No Row(s) Affected")
	}

	return nil
}

func (a *AuthRepo) SetResetToken(ctx context.Context, hashedToken string, userID string, ttl time.Duration) error {
	key := os.Getenv("REDIS_USER") + ":reset_password:" + hashedToken
	return a.rdb.Set(ctx, key, userID, ttl).Err()
}

func (a *AuthRepo) GetResetToken(ctx context.Context, hashedToken string) (string, error) {
	key := os.Getenv("REDIS_USER") + ":reset_password:" + hashedToken
	return a.rdb.Get(ctx, key).Result()
}

func (a *AuthRepo) DeleteResetToken(ctx context.Context, hashedToken string) error {
	key := os.Getenv("REDIS_USER") + ":reset_password:" + hashedToken
	return a.rdb.Del(ctx, key).Err()
}

func (a *AuthRepo) UpdatePassword(ctx context.Context, db DBTX, userID string, newHashedPassword string) error {
	query := "UPDATE users SET password = $1 WHERE id = $2"
	cmd, err := db.Exec(ctx, query, newHashedPassword, userID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("User not found or password not updated")
	}

	return nil
}

func (a *AuthRepo) BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error {
	key := os.Getenv("REDIS_USER") + ":blacklist:" + jti
	return a.rdb.Set(ctx, key, "revoked", ttl).Err()
}