package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/service"
)

type Handler struct {
	paymentService service.PaymentService
}

func NewHandler(paymentService service.PaymentService) *Handler {
	return &Handler{
		paymentService: paymentService,
	}
}

// HealthCheck 健康檢查
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreatePayment 創建支付
func (h *Handler) CreatePayment(c *gin.Context) {
	var req model.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 從 context 中獲取 userID
	userID := c.GetString("userID")
	req.UserID = userID

	payment, err := h.paymentService.CreatePayment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// GetPayment 獲取支付詳情
func (h *Handler) GetPayment(c *gin.Context) {
	id := c.Param("id")
	payment, err := h.paymentService.GetPayment(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrPaymentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// GetPaymentByOrderID 根據訂單ID獲取支付
func (h *Handler) GetPaymentByOrderID(c *gin.Context) {
	orderID := c.Param("orderId")
	payment, err := h.paymentService.GetPaymentByOrderID(c.Request.Context(), orderID)
	if err != nil {
		if err == service.ErrPaymentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// GetUserPayments 獲取用戶支付記錄
func (h *Handler) GetUserPayments(c *gin.Context) {
	userID := c.Param("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	payments, err := h.paymentService.GetUserPayments(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// ListPayments 獲取支付列表
func (h *Handler) ListPayments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	payments, err := h.paymentService.ListPayments(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// ProcessPayment 處理支付
func (h *Handler) ProcessPayment(c *gin.Context) {
	id := c.Param("id")
	err := h.paymentService.ProcessPayment(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrPaymentNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		case service.ErrInvalidPaymentStatus:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment status"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment processed successfully"})
}

// RefundPayment 退款
func (h *Handler) RefundPayment(c *gin.Context) {
	var req model.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.paymentService.RefundPayment(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrPaymentNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		case service.ErrInvalidPaymentStatus:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment status"})
		case service.ErrInvalidRefundAmount:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid refund amount"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "refund processed successfully"})
}

// CancelPayment 取消支付
func (h *Handler) CancelPayment(c *gin.Context) {
	id := c.Param("id")
	err := h.paymentService.CancelPayment(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrPaymentNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		case service.ErrInvalidPaymentStatus:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment status"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment cancelled successfully"})
}
