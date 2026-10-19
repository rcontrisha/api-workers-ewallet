package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-workers-ewallet/internal/controller"
	"github.com/rcontrisha/api-workers-ewallet/internal/middleware"
	"github.com/rcontrisha/api-workers-ewallet/internal/repository"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
	"github.com/redis/go-redis/v9"
)

func AuthRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client, mailer *pkg.Mailer) {
	authRouter := r.Group("auth")

	repo := repository.NewAuthRepo(rdb)
	service := service.NewAuthService(repo, db, mailer)
	controller := controller.NewAuthController(service)

	authRouter.POST("login", controller.LoginController)
	authRouter.POST("register", controller.RegisterController)
	authRouter.POST("forgot-password", controller.ForgotPassword)
	authRouter.POST("reset-password", controller.ResetPassword)

	authRouter.POST("logout", middleware.AuthMiddleware(rdb), controller.Logout)
}
