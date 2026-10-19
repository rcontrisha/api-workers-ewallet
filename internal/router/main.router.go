package router

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/rcontrisha/api-workers-ewallet/docs"
	"github.com/rcontrisha/api-workers-ewallet/internal/middleware"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MainRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client, mailer *pkg.Mailer, dokuClient *pkg.DokuClient) {
	router.Use(middleware.Cors)
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Static("image", path.Join("public", "image"))

	AuthRouter(router, db, rdb, mailer)
	UserRouter(router, db, rdb)
	PaymentRouter(router, db, rdb, dokuClient)
	TransferRouter(router, db, rdb)
	TransactionHistoryRouter(router, db, rdb)
	FavoriteRouter(router, db, rdb)
}
