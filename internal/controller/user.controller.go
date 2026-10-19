package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	msgerror "github.com/rcontrisha/api-workers-ewallet/internal/error"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Update profile information and optionally upload a profile image.
// @Tags User
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param fullname formData string false "Full name"
// @Param phone formData string false "Phone number"
// @Param email formData string false "Email address"
// @Param avatar_url formData file false "JPG, JPEG or PNG image, maximum 2 MiB"
// @Success 200 {object} dto.Response "Profile updated successfully"
// @Failure 400 {object} dto.Response "Invalid form data or avatar"
// @Failure 401 {object} dto.Response "Unauthorized"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /user/update [patch]
func (u *UserController) UpdateUserProfile(ctx *gin.Context) {
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

	var payload dto.UserUpdateRequest
	if err := ctx.ShouldBindWith(&payload, binding.FormMultipart); err != nil {
		log.Println("[UpdateProfile] Bind error:", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var targetFilePath string
	var avatarPtr *string

	if payload.AvatarUrl != nil {
		const maxFileSize = 2 * 1024 * 1024
		if payload.AvatarUrl.Size > maxFileSize {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "file exceeds size limit (2mb)",
			})
			return
		}

		ext := strings.ToLower(filepath.Ext(payload.AvatarUrl.Filename))
		allowedExt := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		if !allowedExt[ext] {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "unsupported format. Please upload in .jpg, .jpeg, or .png format",
			})
			return
		}

		uploadDir := filepath.Join("public", "images", "profile")
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			log.Println("[UpdateProfile] failed to create directory:", err.Error())
			ctx.JSON(http.StatusInternalServerError, dto.Response{
				Success: false,
				Message: "failed processing file",
			})
			return
		}

		filename := fmt.Sprintf("%s_%d%s", userId, time.Now().UnixNano(), ext)
		targetFilePath = filepath.Join(uploadDir, filename)
		if err := ctx.SaveUploadedFile(payload.AvatarUrl, targetFilePath); err != nil {
			log.Println("[UpdateProfile] save upload file err:", err.Error())
			ctx.JSON(http.StatusInternalServerError, dto.Response{
				Success: false,
				Message: "failed to upload file",
			})
			return
		}

		imageUrl := fmt.Sprintf("images/profile/%s", filename)
		avatarPtr = &imageUrl
	}

	account, err := u.service.UpdateProfileInfo(
		ctx.Request.Context(),
		userId,
		payload.Email,
		dto.UserUpdateResponse{
			Id:        userId,
			FullName:  payload.FullName,
			Phone:     payload.Phone,
			AvatarUrl: avatarPtr,
		},
	)
	if err != nil {
		log.Println("[UpdateProfile] service err:", err.Error())
		if targetFilePath != "" {
			if rmErr := os.Remove(targetFilePath); rmErr != nil {
				log.Println("[UpdateProfile] failed to remove file:", rmErr.Error())
			}
		}
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "failed to update profile",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "profile updated successfully",
		Data:    account,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Returns profile information for the authenticated user.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response "Profile retrieved successfully"
// @Failure 401 {object} dto.Response "Missing or invalid bearer token"
// @Failure 404 {object} dto.Response "Profile not found"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /users/profile [get]
func (u *UserController) GetProfile(ctx *gin.Context) {
	claims, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "please login first",
		})
		return
	}
	token, ok := claims.(pkg.JWTClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid bearer token",
		})
		return
	}
	profile, err := u.service.GetProfileService(ctx.Request.Context(), token.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Profile not found.",
			})
			return
		}
		log.Println("error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get profile success.",
		"data": gin.H{
			"profile": profile,
		},
	})
}

// DeleteAvatar godoc
// @Summary Delete avatar
// @Description Delete the authenticated user avatar.
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response "Avatar deleted"
// @Failure 401 {object} dto.Response "Not authenticated"
// @Failure 404 {object} dto.Response "Profile not found"
// @Failure 500 {object} dto.Response "Server error"
// @Router /user/avatar [delete]
func (u *UserController) DeleteAvatar(ctx *gin.Context) {
	claims, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "token doesn't exist",
		})
		return
	}
	token, ok := claims.(pkg.JWTClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "invalid token",
		})
		return
	}

	if err := u.service.DeleteAvatar(ctx.Request.Context(), token.Id); err != nil {
		log.Println("[DeleteAvatar] service err:", err.Error())
		if errors.Is(err, msgerror.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: "profile not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "failed to delete avatar",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "avatar deleted successfully",
	})
}

// ChangePassword godoc
// @Summary Change password
// @Description Change password using old and new passwords.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dto.UserChangePassword true "Request payload"
// @Success 200 {object} dto.Response "Password changed"
// @Failure 400 {object} dto.Response "Invalid request or old password"
// @Failure 401 {object} dto.Response "Not authenticated"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /user/change-password [patch]
func (u *UserController) ChangePassword(ctx *gin.Context) {
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

	var body dto.UserChangePassword

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("[change password] err:", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := u.service.ChangePassword(ctx.Request.Context(), userId, body); err != nil {
		log.Println("[change password] err:", err.Error())
		if errors.Is(err, msgerror.ErrWrongPass) {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "wrong old password",
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "change password success",
	})
}

// ChangePinUser godoc
// @Summary Change transaction PIN
// @Description Change the authenticated user transaction PIN.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dto.UserChangePin true "Request payload"
// @Success 200 {object} dto.Response "PIN changed"
// @Failure 400 {object} dto.Response "Invalid request"
// @Failure 401 {object} dto.Response "Not authenticated"
// @Failure 500 {object} dto.Response "Server error"
// @Router /user/change-pin [patch]
func (u *UserController) ChangePinUser(ctx *gin.Context) {
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
	var body dto.UserChangePin
	if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		log.Println("[change pin] err :", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "invalid request",
		})
		return
	}
	if err := u.service.ChangePin(ctx.Request.Context(), userId, body); err != nil {
		log.Println("[change pin controller] err: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "an error server occured",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "pin changed successfully",
	})

}

// SetPIN godoc
// @Summary Set transaction PIN
// @Description Set a six-digit numeric transaction PIN.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dto.SetPinRequest true "Request payload"
// @Success 200 {object} dto.Response "PIN set"
// @Failure 400 {object} dto.Response "Invalid PIN or request"
// @Failure 401 {object} dto.Response "Not authenticated"
// @Failure 500 {object} dto.Response "Authentication middleware or internal server error"
// @Router /user/pin [post]
func (u *UserController) SetPIN(ctx *gin.Context) {
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

	var payload dto.SetPinRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err := u.service.SetPINService(ctx.Request.Context(), userId, payload); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "PIN set successfully",
		Data:    nil,
	})
}

// UserSummaryDashboard godoc
// @Summary Get dashboard summary
// @Description Retrieve wallet balance, total income and expense for the dashboard.
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response "Summary retrieved"
// @Failure 401 {object} dto.Response "Not authenticated"
// @Failure 500 {object} dto.Response "Server error"
// @Router /user/dashboard-summary [get]
func (u *UserController) UserSummaryDashboard(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "unauthorized",
		})
		return
	}
	claims := token.(pkg.JWTClaims)
	userId := claims.Id
	data, err := u.service.UserDashboardSummary(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "an error server occured",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "dashboard user",
		Data:    data,
	})
}

// DashboardHistoryTransaction godoc
// @Summary Get dashboard transaction history
// @Description Retrieve recent transactions displayed on the dashboard.
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response "History retrieved"
// @Failure 401 {object} dto.Response "Not authenticated"
// @Failure 500 {object} dto.Response "Server error"
// @Router /user/dashboard-history [get]
func (u *UserController) DashboardHistoryTransaction(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "unauthorized",
		})
		return
	}
	claims := token.(pkg.JWTClaims)
	userId := claims.Id

	data, err := u.service.UserDashboardHistoryTransaction(ctx.Request.Context(), userId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "an server error occured",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "history transactions retrieved successfully",
		Data:    data,
	})
}
