package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	msgerror "github.com/rcontrisha/api-workers-ewallet/internal/error"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type FavoriteController struct {
	service *service.FavoriteService
}

func NewFavoriteController(service *service.FavoriteService) *FavoriteController {
	return &FavoriteController{
		service: service,
	}
}

func (f *FavoriteController) Favorite(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "token doesn't exist",
		})
		return
	}
	claims, ok := token.(pkg.JWTClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "invalid token",
		})
		return
	}
	userId := claims.Id
	var req dto.FavoriteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "favorite_user_id is required and must be a valid uuid",
		})
		return
	}
	isFavorite, err := f.service.ToggleFavorite(ctx.Request.Context(), userId, req.FavoriteUserID)
	if err != nil {
		switch {
		case errors.Is(err, msgerror.ErrUserNotFound):
			ctx.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: "user not found",
			})
		default:
			log.Println("[favorite controller] toggle err:", err)
			ctx.JSON(http.StatusInternalServerError, dto.Response{
				Success: false,
				Message: "an error server occured",
			})
		}
		return
	}
	msg := "success remove favorite"
	if isFavorite {
		msg = "success add favorite"
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: msg,
		Data:    gin.H{"is_favorite": isFavorite},
	})

}
