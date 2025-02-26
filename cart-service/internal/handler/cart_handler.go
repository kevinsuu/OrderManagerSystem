package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	userID := c.GetString("userID")

	if err := h.cartService.AddItem(c.Request.Context(), userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}

// ... 實現其他處理器方法 ...
