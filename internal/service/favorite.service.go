package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-workers-ewallet/internal/repository"
)

type FavoriteService struct {
	repo *repository.FavoriteRepo
	db   *pgxpool.Pool
}

func NewFavoriteService(repo *repository.FavoriteRepo, db *pgxpool.Pool) *FavoriteService {
	return &FavoriteService{
		repo: repo,
		db:   db,
	}

}

func (s *FavoriteService) ToggleFavorite(ctx context.Context, userID, favoriteUserID string) (bool, error) {
	if userID == favoriteUserID {
		return false, errors.New("can't favorite yourself")

	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)
	isFav, err := s.repo.IsFavorite(ctx, tx, userID, favoriteUserID)
	if err != nil {
		return false, err
	}
	if isFav {
		if err := s.repo.RemoveFavorite(ctx, tx, userID, favoriteUserID); err != nil {
			return false, err
		}
	} else {
		if err := s.repo.AddFavorite(ctx, tx, userID, favoriteUserID); err != nil {
			return false, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return !isFav, nil
}
