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

func UserRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	userRouter := router.Group("user")
	userRepo := repository.NewUserRepo()
	userService := service.NewUserService(userRepo, db)
	userController := controller.NewUserController(userService)

	userRouter.PATCH("update", middleware.AuthMiddleware(rdb), userController.UpdateUserProfile)
	userRouter.GET("profile", middleware.AuthMiddleware(rdb), userController.GetProfile)
	userRouter.DELETE("avatar", middleware.AuthMiddleware(rdb), userController.DeleteAvatar)
	userRouter.PATCH("change-password", middleware.AuthMiddleware(rdb), userController.ChangePassword)
	userRouter.PATCH("change-pin", middleware.AuthMiddleware(rdb), userController.ChangePinUser)
	userRouter.POST("pin", middleware.AuthMiddleware(rdb), userController.SetPIN)
	userRouter.GET("dashboard-summary", middleware.AuthMiddleware(rdb), userController.UserSummaryDashboard)
	userRouter.GET("dashboard-history", middleware.AuthMiddleware(rdb), userController.DashboardHistoryTransaction)
}
