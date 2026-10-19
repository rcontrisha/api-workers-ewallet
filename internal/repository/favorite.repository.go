package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	msgerror "github.com/rcontrisha/api-workers-ewallet/internal/error"
)

type FavoriteRepo struct{}

func NewFavoriteRepo() *FavoriteRepo {
	return &FavoriteRepo{}
}

func (f *FavoriteRepo) AddFavorite(ctx context.Context, db DBTX, userId, favoriteUserId string) error {
	query := `INSERT INTO favorite (user_id, favorite_user_id, created_at)
					VALUES ($1, $2, NOW())
					ON CONFLICT (user_id, favorite_user_id) DO NOTHING`
	_, err := db.Exec(ctx, query, userId, favoriteUserId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return msgerror.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *FavoriteRepo) RemoveFavorite(ctx context.Context, db DBTX, userID, favoriteUserID string) error {
	const q = `DELETE FROM favorite WHERE user_id = $1 AND favorite_user_id = $2`

	tag, err := db.Exec(ctx, q, userID, favoriteUserID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return msgerror.ErrFavNotFound
	}
	return nil
}

func (r *FavoriteRepo) IsFavorite(ctx context.Context, db DBTX, userID, favoriteUserID string) (bool, error) {
	const q = `
SELECT EXISTS (
  SELECT 1 FROM favorite
  WHERE user_id = $1 AND favorite_user_id = $2
)`

	var exists bool
	if err := db.QueryRow(ctx, q, userID, favoriteUserID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
