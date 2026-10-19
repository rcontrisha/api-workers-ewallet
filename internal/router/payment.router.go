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

func PaymentRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client, dokuClient *pkg.DokuClient) {
	paymentRouter := r.Group("payment")

	paymentRepo := repository.NewPaymentRepo()
	userRepo := repository.NewUserRepo()

	service := service.NewPaymentService(paymentRepo, userRepo, db, dokuClient)
	controller := controller.NewPaymentController(service)

	paymentRouter.POST("topup", middleware.AuthMiddleware(rdb), controller.CreateTopup)
	paymentRouter.GET("topup/:trx_code", middleware.AuthMiddleware(rdb), controller.GetTopupSummary)

	// DOKU webhook (notify)
	paymentRouter.POST("doku/notify", controller.DokuNotification)
	paymentRouter.GET("method", controller.GetPaymentMethodList)
}
