package service

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	msgerror "github.com/rcontrisha/api-workers-ewallet/internal/error"
	"github.com/rcontrisha/api-workers-ewallet/internal/model"
	"github.com/rcontrisha/api-workers-ewallet/internal/repository"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type UserService struct {
	repo *repository.UserRepo
	db   *pgxpool.Pool
}

func NewUserService(repo *repository.UserRepo, db *pgxpool.Pool) *UserService {
	return &UserService{
		repo: repo,
		db:   db,
	}

}

func (u *UserService) UpdateProfileInfo(ctx context.Context, userId, email string, payload dto.UserUpdateResponse) (dto.UserUpdateResponse, error) {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return dto.UserUpdateResponse{}, err
	}

	defer tx.Rollback(ctx)

	data := model.Account{
		Id:       payload.Id,
		FullName: payload.FullName,
		Phone:    payload.Phone,
		Avatar:   payload.AvatarUrl,
	}

	account, err := u.repo.UpdateProfileInfo(ctx, userId, tx, data)
	if err != nil {
		return dto.UserUpdateResponse{}, err
	}

	userEmail, err := u.repo.UpdateUserEmail(ctx, email, userId, tx)
	if err != nil {
		return dto.UserUpdateResponse{}, err
	}

	profile := dto.UserUpdateResponse{
		Id:        account.Id,
		FullName:  account.FullName,
		Phone:     account.Phone,
		AvatarUrl: account.Avatar,
		Email:     userEmail,
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.UserUpdateResponse{}, err
	}
	return profile, nil

}

func (s *UserService) GetProfileService(ctx context.Context, userID string) (dto.UserResponse, error) {
	user, err := s.repo.GetProfile(ctx, s.db, userID)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return dto.UserResponse{
		ID:        user.ID,
		Fullname:  user.Fullname,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdateAt:  user.UpdatedAt,
	}, nil
}

func (u *UserService) DeleteAvatar(ctx context.Context, userID string) error {
	return u.repo.DeleteAvatar(ctx, userID, u.db)
}

func (u *UserService) ChangePassword(ctx context.Context, userId string, payload dto.UserChangePassword) error {
	pass, err := u.repo.GetUserPassword(ctx, userId, u.db)
	if err != nil {
		return err
	}
	if err := pkg.Compare(payload.OldPassword, pass); err != nil {
		return msgerror.ErrWrongPass
	}

	hashConf := pkg.NewRecommendedHashConfig()
	newPass := hashConf.GenHash(payload.NewPassword)

	return u.repo.ChangePassword(ctx, userId, newPass, u.db)
}

func (u *UserService) ChangePin(ctx context.Context, userId string, payload dto.UserChangePin) error {
	hashConf := pkg.NewRecommendedHashConfig()
	pin := hashConf.GenHash(payload.Pin)

	return u.repo.ChangePinUser(ctx, userId, u.db, pin)

}

func (s *UserService) SetPINService(ctx context.Context, userID string, payload dto.SetPinRequest) error {
	acc, err := s.repo.GetProfile(ctx, s.db, userID)
	if err != nil {
		log.Println(err)
		return err
	}

	if acc.Pin != nil && *acc.Pin != "" {
		return errors.New("PIN has already been set. Please use Change PIN feature.")
	}

	hashConf := pkg.NewRecommendedHashConfig()
	hashedPIN := hashConf.GenHash(payload.Pin)

	if err := s.repo.SetUserPIN(ctx, s.db, userID, hashedPIN); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UserDashboardSummary(ctx context.Context, userId string) (dto.DashboardSummaryResponse, error) {
	data, err := s.repo.UserSummaryDashboard(ctx, userId, s.db)
	if err != nil {
		return dto.DashboardSummaryResponse{}, nil
	}

	payload := dto.DashboardSummaryResponse{
		Balance: data.Balance,
		Income:  data.Income,
		Expense: data.Expense,
	}
	return payload, nil

}

func (s *UserService) UserDashboardHistoryTransaction(ctx context.Context, userId string) ([]dto.HistoryTransactionResponse, error) {
	data, err := s.repo.UserHistoryTransaction(ctx, userId, s.db)
	if err != nil {
		return nil, err
	}
	histories := make([]dto.HistoryTransactionResponse, 0)
	for _, v := range data {
		histories = append(histories, dto.HistoryTransactionResponse{
			Name:   v.Name,
			Avatar: v.Avatar,
			Type:   v.Type,
			Amount: v.Amount,
		})
	}
	return histories, nil
}
