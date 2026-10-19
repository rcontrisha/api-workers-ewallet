package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/model"
	"github.com/rcontrisha/api-workers-ewallet/internal/repository"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type TransferService struct {
	repo *repository.TransferRepo
	db   *pgxpool.Pool
}

func NewTransferService(repo *repository.TransferRepo, db *pgxpool.Pool) *TransferService {
	return &TransferService{
		repo: repo,
		db:   db,
	}
}

func (s *TransferService) SearchUsersService(ctx context.Context, query string, currentUserID string) ([]dto.SearchUserResponse, error) {
	var users []model.Account
	var err error

	if query == "" {
		users, err = s.repo.GetRecentAndFavoriteContacts(ctx, s.db, currentUserID)
	} else {
		users, err = s.repo.SearchUsers(ctx, s.db, query, currentUserID)
	}

	if err != nil {
		return nil, err
	}

	res := make([]dto.SearchUserResponse, 0)
	for _, u := range users {
		res = append(res, dto.SearchUserResponse{
			ID:       u.Id,
			FullName: u.FullName,
			Phone:    u.Phone,
			Avatar:   u.Avatar,
		})
	}
	return res, nil
}

func (s *TransferService) TransferMoneyService(ctx context.Context, senderID string, payload dto.TransferRequest) (dto.TransferResponse, error) {
	if senderID == payload.ReceiverID {
		return dto.TransferResponse{}, errors.New("Cannot transfer to your own wallet.")
	}

	hashedPIN, _, err := s.repo.GetUserPINAndProfile(ctx, s.db, senderID)
	if err != nil {
		if err.Error() == "PIN_NOT_SET" {
			return dto.TransferResponse{}, errors.New("Please set your transaction PIN first.")
		}
		return dto.TransferResponse{}, err
	}

	if err := pkg.Compare(payload.Pin, hashedPIN); err != nil {
		return dto.TransferResponse{}, errors.New("Invalid transaction PIN.")
	}

	_, receiverProfile, err := s.repo.GetUserPINAndProfile(ctx, s.db, payload.ReceiverID)
	if err != nil {
		log.Println(err)
		return dto.TransferResponse{}, errors.New("Receiver user not found.")
	}

	trxCode := fmt.Sprintf("TRF-%d", time.Now().UnixNano()/1e6)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return dto.TransferResponse{}, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err.Error() != "tx is closed" {
			log.Println("[TransferService]-Rollback Error:", err)
		}
	}()

	if err := s.repo.DeductWalletBalance(ctx, tx, senderID, payload.Amount); err != nil {
		if err.Error() == "INSUFFICIENT_BALANCE" {
			return dto.TransferResponse{}, errors.New("Insufficient wallet balance.")
		}
		return dto.TransferResponse{}, err
	}

	if err := s.repo.AddWalletBalance(ctx, tx, payload.ReceiverID, payload.Amount); err != nil {
		return dto.TransferResponse{}, fmt.Errorf("failed to topup receiver: %w", err)
	}

	trModel := model.Transaction{
		TrxCode:          trxCode,
		Type:             "transfer",
		Status:           "success",
		Amount:           payload.Amount,
		Tax:              0,
		Discount:         0,
		Total:            payload.Amount,
		SenderWalletID:   &senderID,
		ReceiverWalletID: payload.ReceiverID,
		Note:             payload.Note,
	}

	createdTrx, err := s.repo.CreateTransaction(ctx, tx, trModel)
	if err != nil {
		return dto.TransferResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.TransferResponse{}, err
	}

	receiverName := "User"
	if receiverProfile.FullName != nil {
		receiverName = *receiverProfile.FullName
	}

	return dto.TransferResponse{
		TrxCode:      createdTrx.TrxCode,
		Amount:       createdTrx.Amount,
		Note:         createdTrx.Note,
		Status:       createdTrx.Status,
		ReceiverName: receiverName,
		CreatedAt:    createdTrx.CreatedAt,
	}, nil
}
