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

// GetTransactionHistory godoc
// @Summary Get transaction history
// @Description Retrieve the authenticated user's successful transaction history with pagination and optional search.
// @Tags History
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number, defaults to 1" default(1)
// @Param limit query int false "Number of transactions per page, defaults to 7" default(7)
// @Param search query string false "Search by name or phone number"
// @Success 200 {object} dto.Response "Transaction history retrieved successfully"
// @Failure 401 {object} dto.Response "Unauthorized or invalid token"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /history [get]
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
// @Summary Hide transaction history
// @Description Hide a transaction from the authenticated user's history without deleting the underlying transaction record.
// @Tags History
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction UUID"
// @Success 200 {object} dto.Response "Transaction hidden successfully"
// @Failure 400 {object} dto.Response "Invalid transaction ID"
// @Failure 401 {object} dto.Response "Unauthorized or invalid token"
// @Failure 404 {object} dto.Response "Transaction not found, inaccessible, or already hidden"
// @Failure 500 {object} dto.Response "Internal server error"
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
