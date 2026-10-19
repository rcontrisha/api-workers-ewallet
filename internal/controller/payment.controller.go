package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type PaymentController struct {
	service *service.PaymentService
}

func NewPaymentController(service *service.PaymentService) *PaymentController {
	return &PaymentController{
		service: service,
	}
}

// CreateTopup godoc
// @Summary Create top-up
// @Description Request a top-up checkout with the DOKU payment gateway and transaction PIN.
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dto.CreateTopupRequest true "Request payload"
// @Success 201 {object} dto.Response "Top-up initiated"
// @Failure 400 {object} dto.Response "Invalid input"
// @Failure 401 {object} dto.Response "Not authenticated"
// @Failure 500 {object} dto.Response "Server error"
// @Router /payment/topup [post]
func (p *PaymentController) CreateTopup(ctx *gin.Context) {
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

	var payload dto.CreateTopupRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, dto.Response{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	resp, err := p.service.CreateTopupService(ctx.Request.Context(), userID, payload)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(201, dto.Response{
		Message: "Topup request generated successfully",
		Data:    resp,
	})
}

// DokuNotification godoc
// @Summary Receive DOKU notification
// @Description Webhook endpoint receiving DOKU transaction status notifications.
// @Tags Payment
// @Accept json
// @Produce json
// @Param payload body dto.DokuNotificationPayload true "Request payload"
// @Success 200 {object} dto.Response "Notification processed"
// @Failure 400 {object} dto.Response "Invalid webhook"
// @Failure 500 {object} dto.Response "Server error"
// @Router /payment/doku/notify [post]
func (p *PaymentController) DokuNotification(c *gin.Context) {
	var payload dto.DokuNotificationPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, dto.Response{
			Message: "Invalid webhook payload",
			Data:    nil,
		})
		return
	}

	if err := p.service.ProcessDokuNotificationService(c.Request.Context(), payload); err != nil {
		c.JSON(500, dto.Response{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(200, dto.Response{
		Message: "Notification processed successfully",
		Data:    nil,
	})
}

// GetPaymentMethodList godoc
// @Summary List payment methods
// @Description Retrieve available payment methods.
// @Tags Payment
// @Produce json
// @Success 200 {object} dto.Response "Payment methods retrieved"
// @Failure 500 {object} dto.Response "Server error"
// @Router /payment/method [get]
func (p *PaymentController) GetPaymentMethodList(ctx *gin.Context) {
	paymentMethod, err := p.service.GetPaymentMethodList(ctx)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: "an server error occured",
		})
		return
	}
	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "payment method list retrieved successfully",
		Data:    paymentMethod,
	})
}

// GetTopupSummary godoc
// @Summary Get top-up transaction summary
// @Description Retrieve the top-up transaction summary by transaction code for the authenticated user.
// @Tags Payment
// @Produce json
// @Security BearerAuth
// @Param trx_code path string true "Top-up transaction code"
// @Success 200 {object} dto.Response "Top-up summary retrieved successfully"
// @Failure 401 {object} dto.Response "Unauthorized or invalid token"
// @Failure 404 {object} dto.Response "Transaction not found"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /payment/topup/{trx_code} [get]
func (p *PaymentController) GetTopupSummary(ctx *gin.Context) {
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

	trxCode := ctx.Param("trx_code")
	resp, err := p.service.GetTopupSummaryService(ctx.Request.Context(), claims.Id, trxCode)
	if err != nil {
		if errors.Is(err, service.ErrTransactionNotFound) {
			ctx.JSON(404, dto.Response{Success: false, Message: "transaction not found"})
			return
		}
		ctx.JSON(500, dto.Response{Success: false, Message: "an server error occured"})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "topup summary retrieved successfully",
		Data:    resp,
	})
}
