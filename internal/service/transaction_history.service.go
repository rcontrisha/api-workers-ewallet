package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/repository"
)

var ErrHistoryNotFound = errors.New("transaction history not found")

type TransactionHistoryService struct {
	repo *repository.TransactionHistoryRepo
	db   *pgxpool.Pool
}

func NewTransactionHistoryService(repo *repository.TransactionHistoryRepo, db *pgxpool.Pool) *TransactionHistoryService {
	return &TransactionHistoryService{
		repo: repo,
		db:   db,
	}

}

func (s *TransactionHistoryService) TransactionHistory(ctx context.Context, walletID, endPoint string, qp dto.TransactionHistoryQuery) ([]dto.TransactionHistory, dto.PaginationMetaData, error) {
	qp.Search = strings.TrimSpace(qp.Search)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, dto.PaginationMetaData{}, err
	}
	defer tx.Rollback(ctx)

	offset := (qp.Page - 1) * qp.Limit

	result, err := s.repo.GetHistoryTransaction(ctx, tx, walletID, qp.Search, qp.Limit, offset)
	if err != nil {
		log.Println(err.Error())
		return nil, dto.PaginationMetaData{}, err
	}

	data := make([]dto.TransactionHistory, 0, len(result))
	for _, v := range result {
		data = append(data, dto.TransactionHistory{
			ID:     v.ID,
			Name:   v.Name,
			Avatar: v.Avatar,
			Phone:  v.Phone,
			Type:   v.Type,
			Amount: v.Amount,
		})
	}

	totalData, err := s.repo.GetTotalDataTransactionHistory(ctx, tx, walletID, qp.Search)
	if err != nil {
		log.Println(err.Error())
		return nil, dto.PaginationMetaData{}, err
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(qp.Limit)))

	var next, prev string
	if qp.Page < totalPage {
		next = buildPageURL(endPoint, qp.Search, qp.Limit, qp.Page+1)
	}
	if qp.Page > 1 {
		prev = buildPageURL(endPoint, qp.Search, qp.Limit, qp.Page-1)
	}

	meta := dto.PaginationMetaData{
		CurrentPage: qp.Page,
		TotalPage:   float64(totalPage),
		TotalData:   totalData,
		Next:        next,
		Prev:        prev,
	}

	return data, meta, nil
}

func (s *TransactionHistoryService) DeleteHistoryTransaction(ctx context.Context, transactionID, walletID string) error {
	err := s.repo.DeleteHistoryTransaction(ctx, s.db, transactionID, walletID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrHistoryNotFound
	}
	return err
}

func buildPageURL(endPoint, search string, limit, page int) string {
	q := url.Values{}
	q.Set("limit", strconv.Itoa(limit))
	q.Set("page", strconv.Itoa(page))
	if search != "" {
		q.Set("search", search)
	}
	return fmt.Sprintf("%s?%s", endPoint, q.Encode())
}
