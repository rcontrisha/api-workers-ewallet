package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"uuid"

	"github.com/gin-gonic/gin"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type TransactionHistoryController struct {
	service *service.TransactionHistoryService
}

func NewTransactionHistoryController(service *service.TransactionHistoryService) *TransactionHistoryController {
	return &TransactionHistoryController{
		service: service,
	}

}

func (c *TransactionHistoryController) GetTransactionHistory(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "7"))
	if err != nil || limit < 1 {
		limit = 7
	}
	search := ctx.Query("search")

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
	qp := dto.TransactionHistoryQuery{
		Search: search,
		Limit:  limit,
		Page:   page,
	}

	trxHistories, metaData, err := c.service.TransactionHistory(ctx.Request.Context(), userId, ctx.Request.URL.Path, qp)
	if err != nil {
		log.Println("[transaction history controller] err: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "an error server occured",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "success get transaction histories",
		Data:    trxHistories,
		Meta:    metaData,
	})
}

// DeleteHistoryTransaction godoc
// @Summary Hide a transaction from the authenticated user's history
// @Description Soft-delete a transaction for the current wallet owner only.
// @Tags History
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction UUID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response "Invalid transaction ID"
// @Failure 401 {object} dto.Response "Not authenticated"
// @Failure 404 {object} dto.Response "Transaction not found in user's history"
// @Failure 500 {object} dto.Response "Server error"
// @Router /history/{id} [delete]
func (c *TransactionHistoryController) DeleteHistoryTransaction(ctx *gin.Context) {
	token, exists := ctx.Get("token")
	if !exists {
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

	transactionID := ctx.Param("id")
	if _, err := uuid.Parse(transactionID); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "invalid transaction ID",
		})
		return
	}

	err := c.service.DeleteHistoryTransaction(ctx.Request.Context(), transactionID, claims.Id)
	if err != nil {
		if errors.Is(err, service.ErrHistoryNotFound) {
			ctx.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: "transaction not found in user's history",
			})
			return
		}
		log.Println("[delete transaction history] err:", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "failed to delete transaction history",
		})
		return
	}


	
	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "transaction history deleted successfully",
	})
}
