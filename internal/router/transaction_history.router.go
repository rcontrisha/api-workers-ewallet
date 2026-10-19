package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-workers-ewallet/internal/controller"
	"github.com/rcontrisha/api-workers-ewallet/internal/middleware"
	"github.com/rcontrisha/api-workers-ewallet/internal/repository"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/redis/go-redis/v9"
)

func TransactionHistoryRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	trxHistoryRouter := router.Group("history")
	trxHistoryRepo := repository.NewTransactionHistoryRepo()
	trxHistoryService := service.NewTransactionHistoryService(trxHistoryRepo, db)
	trxHistoryController := controller.NewTransactionHistoryController(trxHistoryService)

	trxHistoryRouter.GET("", middleware.AuthMiddleware(rdb), trxHistoryController.GetTransactionHistory)
	trxHistoryRouter.DELETE(":id", middleware.AuthMiddleware(rdb), trxHistoryController.DeleteHistoryTransaction)
}
