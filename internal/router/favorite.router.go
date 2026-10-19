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

func FavoriteRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	favoriteRouter := router.Group("favorite")
	favoriteRepo := repository.NewFavoriteRepo()
	favoriteSrv := service.NewFavoriteService(favoriteRepo, db)
	favoriteCtrl := controller.NewFavoriteController(favoriteSrv)
	favoriteRouter.POST("", middleware.AuthMiddleware(rdb), favoriteCtrl.Favorite)
}
