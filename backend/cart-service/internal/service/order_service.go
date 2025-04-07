package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/repository"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("invalid order status")
)

// OrderService 訂單服務接口
type OrderService interface {
	CreateOrder(ctx context.Context, userID string, cartItems []model.OrderItem, shippingInfo model.ShippingInfo) (*model.Order, error)
	GetOrder(ctx context.Context, orderID string) (*model.Order, error)
	GetUserOrders(ctx context.Context, userID string, page, limit int) ([]model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status model.OrderStatus) error
}

// orderService 訂單服務實現
type orderService struct {
	orderRepo repository.OrderRepository
}

// NewOrderService 創建新的訂單服務實例
func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

// CreateOrder 創建新訂單
func (s *orderService) CreateOrder(ctx context.Context, userID string, cartItems []model.OrderItem, shippingInfo model.ShippingInfo) (*model.Order, error) {
	// 計算訂單總金額
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.TotalPrice
	}

	order := &model.Order{
		ID:           uuid.New().String(),
		UserID:       userID,
		Items:        cartItems,
		TotalAmount:  totalAmount,
		Status:       model.OrderStatusPending,
		ShippingInfo: shippingInfo,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder 獲取訂單詳情
func (s *orderService) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

// GetUserOrders 獲取用戶的所有訂單
func (s *orderService) GetUserOrders(ctx context.Context, userID string, page, limit int) ([]model.Order, error) {
	// 計算偏移量
	offset := (page - 1) * limit
	return s.orderRepo.GetByUserID(ctx, userID, offset, limit)
}

// UpdateOrderStatus 更新訂單狀態
func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID string, status model.OrderStatus) error {
	// 驗證訂單狀態是否有效
	switch status {
	case model.OrderStatusPending,
		model.OrderStatusPaid,
		model.OrderStatusShipped,
		model.OrderStatusDelivered,
		model.OrderStatusCancelled:
		// 有效狀態
	default:
		return ErrInvalidOrderStatus
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, status)
}
