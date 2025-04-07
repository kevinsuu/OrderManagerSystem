package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/service"
)

type Handler struct {
	authService service.IAuthService
}

func NewHandler(authService service.IAuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

type ForgotPasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

// ResetPasswordRequest 定義重設密碼請求結構
type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

// ForgotPassword 處理忘記密碼請求
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式"})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "新密碼與確認密碼不符"})
		return
	}

	// 調用服務層重設密碼
	err := h.authService.ForgetPassword(c.Request.Context(), req.Email, req.NewPassword)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "找不到此電子郵件地址的帳號"})
		case service.ErrInvalidPassword:
			c.JSON(http.StatusBadRequest, gin.H{"error": "密碼格式不正確"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "重設密碼失敗，請稍後再試"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "密碼重設成功",
	})
}

// ResetPassword 處理重設密碼請求
func (h *Handler) ResetPassword(c *gin.Context) {
	// 從 Authorization header 獲取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要登入才能重設密碼"})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的認證格式"})
		return
	}

	var req struct {
		NewPassword     string `json:"newPassword" binding:"required,min=6"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式"})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "新密碼與確認密碼不符"})
		return
	}

	// 調用服務層重設密碼
	err := h.authService.ResetPassword(c.Request.Context(), tokenParts[1], req.NewPassword)
	if err != nil {
		switch err {
		case service.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的認證令牌"})
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "找不到該用戶"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "重設密碼失敗"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "密碼重設成功",
	})
}

// Register 用戶註冊
func (h *Handler) Register(c *gin.Context) {
	var req model.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login 用戶登錄
func (h *Handler) Login(c *gin.Context) {
	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case service.ErrInvalidPassword:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// ValidateToken 驗證令牌
func (h *Handler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		return
	}

	user, err := h.authService.ValidateToken(c.Request.Context(), tokenParts[1])
	if err != nil {
		switch err {
		case service.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		case service.ErrTokenExpired:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// RefreshToken 刷新令牌
func (h *Handler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token is required"})
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case service.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUser 獲取用戶信息
func (h *Handler) GetUser(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateAddress 添加新的處理方法
func (h *Handler) CreateAddress(c *gin.Context) {
	var req model.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	address, err := h.authService.CreateAddress(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, address)
}

// DeleteAddress 刪除地址
func (h *Handler) DeleteAddress(c *gin.Context) {
	userID := c.GetString("userID")
	addressID := c.Param("id")

	// 檢查地址是否存在並屬於該用戶
	address, err := h.authService.GetAddressByID(c.Request.Context(), addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if address == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
		return
	}
	if address.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// 執行刪除操作
	if err := h.authService.DeleteAddress(c.Request.Context(), userID, addressID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "address deleted successfully"})
}

// GetAddresses 獲取用戶地址列表
func (h *Handler) GetAddresses(c *gin.Context) {
	userID := c.GetString("userID")

	addresses, err := h.authService.GetAddresses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// UpdateAddress 更新地址
func (h *Handler) UpdateAddress(c *gin.Context) {
	var req model.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	addressID := c.Param("id")

	address, err := h.authService.UpdateAddress(c.Request.Context(), userID, addressID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, address)
}

// GetPreference 獲取用戶偏好設置
func (h *Handler) GetPreference(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	pref, err := h.authService.GetPreference(c.Request.Context(), userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pref)
}

// UpdatePreference 更新用戶偏好設置
func (h *Handler) UpdatePreference(c *gin.Context) {
	var req model.PreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	pref, err := h.authService.UpdatePreference(c.Request.Context(), userID, &req)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pref)
}

// SetDefaultAddress 設置預設地址
func (h *Handler) SetDefaultAddress(c *gin.Context) {
	userID := c.GetString("userID")
	addressID := c.Param("id")

	address, err := h.authService.SetDefaultAddress(c.Request.Context(), userID, addressID)
	if err != nil {
		if err.Error() == "address not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "地址不存在"})
			return
		}
		if err.Error() == "address does not belong to the user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "沒有權限修改此地址"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "已成功設置為預設地址",
		"address": address,
	})
}
