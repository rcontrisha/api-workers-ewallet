package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/model"
	"github.com/rcontrisha/api-workers-ewallet/internal/repository"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type PaymentService struct {
	repo       *repository.PaymentRepo
	userRepo   *repository.UserRepo
	db         *pgxpool.Pool
	dokuClient *pkg.DokuClient
}

var ErrTransactionNotFound = errors.New("transaction not found")

func NewPaymentService(repo *repository.PaymentRepo, userRepo *repository.UserRepo, db *pgxpool.Pool, dokuClient *pkg.DokuClient) *PaymentService {
	return &PaymentService{
		repo:       repo,
		userRepo:   userRepo,
		db:         db,
		dokuClient: dokuClient,
	}
}

func (s *PaymentService) CreateTopupService(ctx context.Context, userID string, payload dto.CreateTopupRequest) (dto.CreateTopupResponse, error) {
	account, err := s.userRepo.GetProfile(ctx, s.db, userID)
	if err != nil {
		return dto.CreateTopupResponse{}, fmt.Errorf("user account not found: %w", err)
	}

	if account.Pin == nil || *account.Pin == "" {
		return dto.CreateTopupResponse{}, fmt.Errorf("user has not set a PIN yet")
	}

	if err := pkg.Compare(payload.Pin, *account.Pin); err != nil {
		log.Printf("[*account.Pin]: %s\n[payload.Pin]: %s", *account.Pin, payload.Pin)
		log.Println(err)
		return dto.CreateTopupResponse{}, fmt.Errorf("invalid PIN")
	}

	var paymentChannels []string
	if payload.PaymentMethodID != nil && *payload.PaymentMethodID != "" {
		pm, err := s.repo.GetPaymentMethodByID(ctx, s.db, *payload.PaymentMethodID)
		if err == nil && pm.IsActive {
			paymentChannels = append(paymentChannels, pm.GatewayPaymentChannel)
		}
	}

	trxCode := fmt.Sprintf("topup-%d", time.Now().UnixNano()/1e6)

	txModel := model.Transaction{
		TrxCode:          trxCode,
		Type:             "topup",
		Status:           "PENDING",
		Amount:           payload.Amount,
		Tax:              payload.Tax,
		Discount:         0,
		Total:            payload.Amount + payload.Tax,
		SenderWalletID:   nil,
		ReceiverWalletID: userID,
		PaymentMethodID:  payload.PaymentMethodID,
		Note:             payload.Note,
	}

	createdTrx, err := s.repo.CreateTransaction(ctx, s.db, txModel)
	if err != nil {
		return dto.CreateTopupResponse{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	dokuPayload := pkg.DokuCheckoutPayload{}
	dokuPayload.Order.Amount = txModel.Total
	dokuPayload.Order.InvoiceNumber = trxCode
	dokuPayload.Order.Currency = "IDR"
	dokuPayload.Order.AutoRedirect = true
	dokuPayload.Order.CallbackURL = fmt.Sprintf("%s/top-up/summary?tranx_id=%s", os.Getenv("DOKU_CALLBACK_URL"), trxCode)

	dokuPayload.Payment.PaymentDueDate = 60
	if len(paymentChannels) > 0 {
		dokuPayload.Payment.PaymentMethodTypes = paymentChannels
	}

	fullName := "Customer"
	if account.Fullname != nil {
		fullName = *account.Fullname
	}
	phone := "081234567890"
	if account.Phone != nil {
		phone = *account.Phone
	}

	dokuPayload.Customer.ID = userID
	dokuPayload.Customer.Name = fullName
	dokuPayload.Customer.Email = account.Email
	dokuPayload.Customer.Phone = phone

	dokuResp, err := s.dokuClient.CreateCheckout(ctx, dokuPayload)
	if err != nil {
		log.Println("[DokuCheckout]-Error: ", err)
		return dto.CreateTopupResponse{}, fmt.Errorf("failed to generate payment checkout: %w", err)
	}

	return dto.CreateTopupResponse{
		TrxCode:    createdTrx.TrxCode,
		Amount:     createdTrx.Amount,
		Total:      createdTrx.Total,
		Status:     createdTrx.Status,
		PaymentURL: dokuResp.Response.Payment.URL,
	}, nil
}

func (s *PaymentService) ProcessDokuNotificationService(ctx context.Context, payload dto.DokuNotificationPayload) error {
	trx, err := s.repo.FindTransactionByTrxCode(ctx, s.db, payload.Order.InvoiceNumber)
	if err != nil {
		return fmt.Errorf("transaction not found: %w", err)
	}

	// Idempotency check: jika status transaksi sudah bukan PENDING, abaikan request duplikat
	if trx.Status != "PENDING" {
		log.Printf("[DokuNotification] Invoice %s already processed with status %s", trx.TrxCode, trx.Status)
		return nil
	}

	if payload.Transaction.Status == "success" {
		tx, err := s.db.Begin(ctx)
		if err != nil {
			return err
		}
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		if err := s.repo.UpdateTransactionStatus(ctx, tx, trx.TrxCode, "success"); err != nil {
			return err
		}

		if err := s.repo.AddWalletBalance(ctx, tx, trx.ReceiverWalletID, trx.Amount); err != nil {
			return err
		}

		if err := tx.Commit(ctx); err != nil {
			return err
		}

		log.Printf("[DokuNotification] Success Topup Invoice %s amount %d for wallet %s", trx.TrxCode, trx.Amount, trx.ReceiverWalletID)
		return nil
	}

	return s.repo.UpdateTransactionStatus(ctx, s.db, trx.TrxCode, "FAILED")
}

func (s *PaymentService) GetPaymentMethodList(ctx context.Context) ([]dto.PaymentMethod, error) {
	data, err := s.repo.GetPaymentMethodList(ctx, s.db)
	if err != nil {
		return nil, err
	}
	methodList := make([]dto.PaymentMethod, 0)
	for _, v := range data {
		methodList = append(methodList, dto.PaymentMethod{
			ID:                    v.ID,
			GatewayPaymentChannel: v.GatewayPaymentChannel,
			Name:                  v.Name,
			Logo:                  v.Logo,
			IsActive:              v.IsActive,
		})
	}
	return methodList, nil
}

func (s *PaymentService) GetTopupSummaryService(ctx context.Context, userID, trxCode string) (dto.TopupSummaryResponse, error) {
	trx, err := s.repo.FindTransactionByTrxCode(ctx, s.db, trxCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.TopupSummaryResponse{}, ErrTransactionNotFound
		}
		return dto.TopupSummaryResponse{}, err
	}

	if trx.Type != "topup" || trx.ReceiverWalletID != userID {
		return dto.TopupSummaryResponse{}, ErrTransactionNotFound
	}

	resp := dto.TopupSummaryResponse{
		TrxCode:   trx.TrxCode,
		Status:    trx.Status,
		Amount:    trx.Amount,
		Tax:       trx.Tax,
		Total:     trx.Total,
		Note:      trx.Note,
		CreatedAt: trx.CreatedAt,
	}

	if trx.PaymentMethodID != nil && *trx.PaymentMethodID != "" {
		if pm, err := s.repo.GetPaymentMethodByID(ctx, s.db, *trx.PaymentMethodID); err == nil {
			resp.PaymentMethod = &pm.Name
			resp.PaymentLogo = pm.Logo
		}
	}

	if account, err := s.userRepo.GetProfile(ctx, s.db, userID); err == nil {
		resp.DestinationPhone = account.Phone
	}

	return resp, nil
}
