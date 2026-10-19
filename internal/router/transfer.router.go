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

func TransferRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	transferGroup := r.Group("transfer", middleware.AuthMiddleware(rdb))

	repo := repository.NewTransferRepo()
	service := service.NewTransferService(repo, db)
	controller := controller.NewTransferController(service)

	transferGroup.GET("search", controller.SearchUsers)
	transferGroup.POST("", controller.TransferMoney)
}
