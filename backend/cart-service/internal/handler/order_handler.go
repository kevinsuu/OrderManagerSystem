package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/service"
)

// OrderHandler 訂單處理器
type OrderHandler struct {
	orderService service.OrderService
	cartService  service.CartService
}

// NewOrderHandler 創建新的訂單處理器
func NewOrderHandler(orderService service.OrderService, cartService service.CartService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		cartService:  cartService,
	}
}

// CreateOrderRequest 創建訂單請求
type CreateOrderRequest struct {
	ShippingInfo model.ShippingInfo `json:"shippingInfo"`
}

// CreateOrder 創建訂單
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetString("userID")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 獲取用戶的購物車
	cartItems, err := h.cartService.GetCartItems(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart items"})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// 將購物車項目轉換為訂單項目
	var orderItems []model.OrderItem
	for _, item := range cartItems {
		orderItem := model.OrderItem{
			ProductID:  item.ProductID,
			Name:       item.Name,
			Price:      item.Price,
			Quantity:   item.Quantity,
			TotalPrice: item.Price * float64(item.Quantity),
		}
		orderItems = append(orderItems, orderItem)
	}

	// 創建訂單
	order, err := h.orderService.CreateOrder(c, userID, orderItems, req.ShippingInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// 清空購物車
	if err := h.cartService.ClearCart(c, userID); err != nil {
		// 記錄錯誤但不影響訂單創建
		// TODO: 添加日誌記錄
	}

	c.JSON(http.StatusOK, order)
}

// GetOrder 獲取訂單詳情
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID := c.GetString("userID")
	orderID := c.Param("id")

	order, err := h.orderService.GetOrder(c, orderID)
	if err != nil {
		if err == service.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order"})
		return
	}

	// 驗證訂單所有者
	if order.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetUserOrders 獲取用戶的所有訂單
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := c.GetString("userID")

	// 獲取分頁參數
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = limitNum
		}
	}

	orders, err := h.orderService.GetUserOrders(c, userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"page":   page,
		"limit":  limit,
	})
}

// UpdateOrderStatus 更新訂單狀態
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req struct {
		Status model.OrderStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.orderService.UpdateOrderStatus(c, orderID, req.Status); err != nil {
		if err == service.ErrInvalidOrderStatus {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status"})
			return
		}
		if err == service.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// ListOrders 獲取訂單列表
func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID := c.GetString("userID")

	// 獲取分頁參數
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = limitNum
		}
	}

	orders, err := h.orderService.GetUserOrders(c, userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"page":   page,
		"limit":  limit,
	})
}

// UpdateOrder 更新訂單
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	orderID := c.Param("id")
	var req struct {
		Status model.OrderStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.orderService.UpdateOrderStatus(c, orderID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

// DeleteOrder 刪除訂單
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Delete operation not supported"})
}

// CancelOrder 取消訂單
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	if err := h.orderService.UpdateOrderStatus(c, orderID, model.OrderStatusCancelled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

// GetOrdersByStatus 根據狀態獲取訂單
func (h *OrderHandler) GetOrdersByStatus(c *gin.Context) {
	status := model.OrderStatus(c.Param("status"))
	userID := c.GetString("userID")

	// 獲取分頁參數
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = limitNum
		}
	}

	orders, err := h.orderService.GetUserOrders(c, userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	// 過濾指定狀態的訂單
	var filteredOrders []model.Order
	for _, order := range orders {
		if order.Status == status {
			filteredOrders = append(filteredOrders, order)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": filteredOrders,
		"page":   page,
		"limit":  limit,
	})
}
