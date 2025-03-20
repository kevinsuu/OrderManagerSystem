package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/client"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/service"
)

type Handler struct {
	cartService service.CartService
}

func NewHandler(cartService service.CartService) *Handler {
	return &Handler{
		cartService: cartService,
	}
}

// GetCart 獲取購物車
func (h *Handler) GetCart(c *gin.Context) {
	userID := c.GetString("userID") // 從 JWT token 中獲取

	cart, err := h.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddToCart 添加商品到購物車
func (h *Handler) AddToCart(c *gin.Context) {
	var req model.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 獲取 token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
		return
	}

	// 獲取用戶ID
	userID := c.GetString("userID")
	log.Printf("AddToCart handler called for user: %s, product: %s", userID, req.ProductID)

	if userID == "" {
		log.Printf("Error: userID is empty in AddToCart handler")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in token"})
		return
	}

	// 直接將完整的 Authorization header 傳遞給 context
	ctx := context.WithValue(c.Request.Context(), client.TokenKey, token)

	if err := h.cartService.AddItem(ctx, userID, &req); err != nil {
		log.Printf("Error in AddItem service: %v", err)
		switch {
		case err == service.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		case strings.Contains(err.Error(), "insufficient stock"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "total quantity exceeds stock"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "unauthorized"):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("Item successfully added to cart for user: %s", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}

// RemoveFromCart 從購物車中移除商品
func (h *Handler) RemoveFromCart(c *gin.Context) {
	userID := c.GetString("userID")
	productID := c.Param("productId")

	if err := h.cartService.RemoveItem(c.Request.Context(), userID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}

// UpdateQuantity 更新購物車商品數量
func (h *Handler) UpdateQuantity(c *gin.Context) {
	var req model.UpdateQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 獲取 token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
		return
	}

	// 直接將完整的 Authorization header 傳遞給 context
	ctx := context.WithValue(c.Request.Context(), client.TokenKey, token)
	userID := c.GetString("userID")

	if err := h.cartService.UpdateQuantity(ctx, userID, &req); err != nil {
		switch {
		case err == service.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		case err == service.ErrInvalidStock:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "unauthorized"):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item quantity updated successfully"})
}

// SelectItems 選擇購物車中的商品
func (h *Handler) SelectItems(c *gin.Context) {
	var req model.SelectItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
		return
	}

	// 直接將完整的 Authorization header 傳遞給 context
	ctx := context.WithValue(c.Request.Context(), client.TokenKey, token)
	userID := c.GetString("userID")

	if err := h.cartService.SelectItems(ctx, userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart items selected successfully"})
}

// ClearCart 清空購物車
func (h *Handler) ClearCart(c *gin.Context) {
	userID := c.GetString("userID")

	if err := h.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}

// CreateOrder 從購物車創建訂單
func (h *Handler) CreateOrder(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), client.TokenKey, token)
	userID := c.GetString("userID")

	if err := h.cartService.CreateOrder(ctx, userID); err != nil {
		switch {
		case strings.Contains(err.Error(), "no items selected"):
			c.JSON(http.StatusBadRequest, gin.H{"error": "no items selected in cart"})
		case strings.Contains(err.Error(), "unauthorized"):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully"})
}

// ... 實現其他處理器方法 ...
