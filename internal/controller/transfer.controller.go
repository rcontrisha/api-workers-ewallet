package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type TransferController struct {
	service *service.TransferService
}

func NewTransferController(service *service.TransferService) *TransferController {
	return &TransferController{
		service: service,
	}
}

// SearchUsers godoc
// @Summary Search transfer recipients
// @Description Search recipient accounts by q. If empty, return recent and favorite contacts.
// @Tags Transfer
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search term (optional)"
// @Success 200 {object} dto.Response "Recipients retrieved"
// @Failure 401 {object} dto.Response "Not authenticated"
// @Failure 500 {object} dto.Response "Server error"
// @Router /transfer/search [get]
func (t *TransferController) SearchUsers(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "token doesn't exist",
		})
		return
	}
	claims, ok := token.(pkg.JWTClaims)
	if !ok {
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "invalid token",
		})
		return
	}
	userID := claims.Id

	query := ctx.Query("q")
	result, err := t.service.SearchUsersService(ctx.Request.Context(), query, userID)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "Users fetched successfully",
		Data:    result,
	})
}

// TransferMoney godoc
// @Summary Transfer money
// @Description Transfer wallet balance to another user after validating the six-digit PIN.
// @Tags Transfer
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dto.TransferRequest true "Request payload"
// @Success 200 {object} dto.Response "Transfer completed"
// @Failure 400 {object} dto.Response "Validation or transfer failed"
// @Failure 401 {object} dto.Response "Not authenticated"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /transfer [post]
func (t *TransferController) TransferMoney(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "token doesn't exist",
		})
		return
	}
	claims, ok := token.(pkg.JWTClaims)
	if !ok {
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "invalid token",
		})
		return
	}
	userID := claims.Id

	var payload dto.TransferRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, dto.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	resp, err := t.service.TransferMoneyService(ctx.Request.Context(), userID, payload)
	if err != nil {
		ctx.JSON(400, dto.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "Transfer completed successfully",
		Data:    resp,
	})
}
