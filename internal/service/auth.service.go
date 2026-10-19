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
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	repo   *repository.AuthRepo
	db     *pgxpool.Pool
	mailer *pkg.Mailer
}

func NewAuthService(repo *repository.AuthRepo, db *pgxpool.Pool, mailer *pkg.Mailer) *AuthService {
	return &AuthService{
		repo:   repo,
		db:     db,
		mailer: mailer,
	}
}

func (a *AuthService) LoginService(ctx context.Context, payload dto.LoginRequest) (dto.LoginResponse, string, error) {
	result, err := a.repo.FindAccount(ctx, a.db, payload.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.LoginResponse{}, "", errors.New("Incorrect Email or Password.	")
		}
		return dto.LoginResponse{}, "", err
	}

	user := dto.LoginResponse{
		Id:     result.Id,
		Email:  result.Email,
		HasPin: result.HasPin,
	}

	if err := pkg.Compare(payload.Password, result.Password); err != nil {
		log.Println(err)
		return dto.LoginResponse{}, "", err
	}

	claims := pkg.NewJWTClaims(result.Id, result.Email)
	token, err := claims.GenToken()

	return user, token, nil
}

func (a *AuthService) RegisterService(ctx context.Context, payload dto.LoginRequest) error {
	_, err := a.repo.FindAccount(ctx, a.db, payload.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if err == nil {
		return errors.New("User Already Exist.")
	}

	hashedPwd := pkg.NewRecommendedHashConfig().GenHash(payload.Password)

	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err.Error())
		}
	}()

	var user = model.User{
		Email:    payload.Email,
		Password: hashedPwd,
	}

	userId, err := a.repo.CreateUser(ctx, tx, user)
	if err != nil {
		log.Println("[Service-CreateUser]-Error: ", err)
		return err
	}
	if err := a.repo.CreateAccount(ctx, tx, userId); err != nil {
		log.Println("[Service-CreateAccount]-Error: ", err)
		return err
	}
	if err := a.repo.CreateWallet(ctx, tx, userId); err != nil {
		log.Println("[Service-CreateWallet]-Error: ", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (a *AuthService) ForgotPasswordService(ctx context.Context, email string) error {
	acc, err := a.repo.FindAccount(ctx, a.db, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Tetap return nil tanpa error biar nggak kena user enumeration attack
			return nil
		}
		return err
	}

	rawToken, err := pkg.GenerateRandomToken(32)
	if err != nil {
		return err
	}

	hashedToken := pkg.HashToken(rawToken)

	if err := a.repo.SetResetToken(ctx, hashedToken, acc.Id, 5*time.Minute); err != nil {
		return err
	}

	go func(toEmail, token string) {
		resetLink := fmt.Sprintf("%s/auth/forgot-password?token=%s", os.Getenv("FE_URL"), token)
		bodyHTML := fmt.Sprintf(`
			<h3>Reset Password Request</h3>
			<p>Silakan klik link di bawah ini untuk mereset password kamu:</p>
			<p><a href="%s">%s</a></p>
			<p>Link ini berlaku selama 5 menit.</p>
		`, resetLink, resetLink)

		if err := a.mailer.Send(toEmail, "Reset Password Request", bodyHTML); err != nil {
			log.Println("[ForgotPasswordMailer]-Error: ", err)
		}
	}(acc.Email, rawToken)

	return nil
}

func (a *AuthService) ResetPasswordService(ctx context.Context, rawToken string, newPassword string) error {
	hashedToken := pkg.HashToken(rawToken)

	userID, err := a.repo.GetResetToken(ctx, hashedToken)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("Token reset password tidak valid atau sudah kadaluarsa.")
		}
		return err
	}

	hashedPwd := pkg.NewRecommendedHashConfig().GenHash(newPassword)

	if err := a.repo.UpdatePassword(ctx, a.db, userID, hashedPwd); err != nil {
		return err
	}

	// Langsung hapus key dari Redis setelah password berhasil diperbarui
	if err := a.repo.DeleteResetToken(ctx, hashedToken); err != nil {
		log.Println("[DeleteResetToken]-Error: ", err)
	}

	return nil
}

func (a *AuthService) Logout(ctx context.Context, jti string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}

	return a.repo.BlacklistToken(ctx, jti, ttl)
}
