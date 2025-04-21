package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/service"
)

// WishlistHandler 處理收藏清單相關請求
type WishlistHandler struct {
	wishlistService service.WishlistService
}

// NewWishlistHandler 創建一個新的收藏清單處理器
func NewWishlistHandler(wishlistService service.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}

// AddToWishlist 添加商品到收藏清單
// @Summary 添加商品到收藏清單
// @Description 將指定商品添加到用戶的收藏清單
// @Tags wishlist
// @Accept json
// @Produce json
// @Param product body model.AddToWishlistRequest true "商品資訊"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/wishlist [post]
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	userId := c.GetString("userID")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授權",
		})
		return
	}

	var req model.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "無效的請求參數",
		})
		return
	}

	if err := h.wishlistService.AddToWishlist(c.Request.Context(), userId, req.ProductId); err != nil {
		// 檢查是否是"已存在於收藏清單"的錯誤
		if err.Error() == "product already in wishlist" {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "商品已在收藏清單中",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("添加到收藏清單失敗: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "已成功添加到收藏清單",
		"data": map[string]bool{
			"alreadyExists": false,
		},
	})
}

// RemoveFromWishlist 從收藏清單移除商品
// @Summary 從收藏清單移除商品
// @Description 從用戶的收藏清單中移除指定商品
// @Tags wishlist
// @Accept json
// @Produce json
// @Param productId path string true "商品ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/wishlist/{productId} [delete]
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	userId := c.GetString("userID")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授權",
		})
		return
	}

	productId := c.Param("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "商品ID不能為空",
		})
		return
	}

	if err := h.wishlistService.RemoveFromWishlist(c.Request.Context(), userId, productId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("從收藏清單移除失敗: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "已成功從收藏清單移除",
	})
}

// GetWishlist 獲取收藏清單
// @Summary 獲取用戶的收藏清單
// @Description 獲取指定用戶的收藏清單
// @Tags wishlist
// @Accept json
// @Produce json
// @Param page query int false "頁碼"
// @Param limit query int false "每頁數量"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/wishlist [get]
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	userId := c.GetString("userID")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授權",
		})
		return
	}

	// 解析分頁參數
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 獲取收藏清單
	wishlistResp, err := h.wishlistService.GetWishlistWithProductDetails(c.Request.Context(), userId, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("獲取收藏清單失敗: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    wishlistResp,
	})
}
