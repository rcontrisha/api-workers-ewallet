package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{
		service: *service,
	}
}

// LoginController godoc
// @Summary Log in
// @Description Authenticate using an email and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body dto.LoginRequest true "Request payload"
// @Success 200 {object} dto.Response "Login succeeded"
// @Failure 401 {object} dto.Response "Invalid email or password"
// @Failure 500 {object} dto.Response "Request parsing or server error"
// @Router /auth/login [post]
func (a *AuthController) LoginController(ctx *gin.Context) {
	var payload dto.LoginRequest
	if e := ctx.ShouldBindWith(&payload, binding.JSON); e != nil {
		log.Println("error", e.Error())
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: "Internal Server Error",
		})
		return
	}

	user, token, err := a.service.LoginService(ctx.Request.Context(), payload)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "Login Failed. Invalid email or password.",
			Data:    gin.H{},
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "Login Success.",
		Data: gin.H{
			"token": token,
			"user":  user,
		},
	})
}

// RegisterController godoc
// @Summary Register an account
// @Description Create a user, account and wallet with email and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body dto.LoginRequest true "Request payload"
// @Success 201 {object} dto.Response "Account created"
// @Failure 401 {object} dto.Response "Registration rejected"
// @Failure 500 {object} dto.Response "Request parsing or server error"
// @Router /auth/register [post]
func (a *AuthController) RegisterController(ctx *gin.Context) {
	var payload dto.LoginRequest

	if e := ctx.ShouldBindWith(&payload, binding.JSON); e != nil {
		log.Println("error", e.Error())
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: "Internal Server Error",
		})
		return
	}

	err := a.service.RegisterService(ctx.Request.Context(), payload)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "Register failed.",
			Data:    gin.H{},
		})
		return
	}

	ctx.JSON(201, dto.Response{
		Success: true,
		Message: "Register Success.",
	})
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset instructions to a registered email.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body dto.ForgotPasswordRequest true "Request payload"
// @Success 200 {object} dto.Response "Instructions sent"
// @Failure 400 {object} dto.Response "Invalid email"
// @Failure 500 {object} dto.Response "Server error"
// @Router /auth/forgot-password [post]
func (a *AuthController) ForgotPassword(c *gin.Context) {
	var payload dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, dto.Response{
			Message: "Invalid email format",
			Data:    nil,
		})
		return
	}

	if err := a.service.ForgotPasswordService(c.Request.Context(), payload.Email); err != nil {
		c.JSON(500, dto.Response{
			Message: "Failed to process forgot password request",
			Data:    nil,
		})
		return
	}

	c.JSON(200, dto.Response{
		Success: true,
		Message: "We have sent password reset instructions to your email.",
		Data:    nil,
	})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using a reset token and a new password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body dto.ResetPasswordRequest true "Request payload"
// @Success 200 {object} dto.Response "Password reset"
// @Failure 400 {object} dto.Response "Invalid or expired reset request"
// @Router /auth/reset-password [post]
func (a *AuthController) ResetPassword(c *gin.Context) {
	var payload dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, dto.Response{
			Message: "Invalid input request. Password must be at least 8 characters.",
			Data:    nil,
		})
		return
	}

	if err := a.service.ResetPasswordService(c.Request.Context(), payload.Token, payload.NewPassword); err != nil {
		c.JSON(400, dto.Response{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(200, dto.Response{
		Message: "Password has been successfully reset. Please login with your new password.",
		Data:    nil,
	})
}

// Logout godoc
// @Summary Log out
// @Description Revoke the authenticated access token using Redis.
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response "Logged out"
// @Failure 401 {object} dto.Response "Not authenticated"
// @Failure 500 {object} dto.Response "Server error"
// @Router /auth/logout [post]
func (a *AuthController) Logout(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "Token Data Not Found in Context.",
		})
	}

	claims := token.(pkg.JWTClaims)
	jti := claims.ID
	expiresAt := claims.ExpiresAt
	if err := a.service.Logout(ctx, jti, expiresAt.Time); err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "Successfully logged out",
	})
}
