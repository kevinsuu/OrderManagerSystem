package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevinsuu/OrderManagerSystem/order-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/order-service/internal/repository"
)

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidOrderState = errors.New("invalid order state")
)

// OrderService 訂單服務接口
type OrderService interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.Order, error)
	GetOrder(ctx context.Context, id string) (*model.OrderResponse, error)
	DeleteOrder(ctx context.Context, id string) error
	ListOrders(ctx context.Context, page, limit int) (*model.OrderListResponse, error)
	GetOrdersByStatus(ctx context.Context, status model.OrderStatus, page, limit int) (*model.OrderListResponse, error)
	CancelOrder(ctx context.Context, id string) error
	GetOrders(ctx context.Context, page, limit int) ([]*model.Order, error)
	GetUserOrders(ctx context.Context, userID string, page, limit int) ([]*model.Order, error)
	UpdateOrder(ctx context.Context, id string, req *model.UpdateOrderRequest) error
	CreateOrderFromCart(ctx context.Context, userID string, req *model.CreateOrderFromCartRequest) (*model.CreateOrderResponse, error)
}

type orderService struct {
	repo  repository.OrderRepository
	redis *repository.RedisRepository
}

// NewOrderService 創建訂單服務實例
func NewOrderService(repo repository.OrderRepository, redis *repository.RedisRepository) OrderService {
	return &orderService{
		repo:  repo,
		redis: redis,
	}
}

// CreateOrder 創建訂單
func (s *orderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.Order, error) {
	order := &model.Order{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Status:    model.OrderStatusPending,
		Items:     req.Items,
		Address:   req.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 計算總金額
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}
	order.TotalAmount = totalAmount

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// 清除相關快取
	s.clearOrderCache(ctx, order.UserID)

	return order, nil
}

// GetOrder 獲取訂單
func (s *orderService) GetOrder(ctx context.Context, id string) (*model.OrderResponse, error) {
	// 嘗試從快取獲取
	cacheKey := fmt.Sprintf("order:%s", id)
	if cached, err := s.redis.Get(ctx, cacheKey); err == nil {
		var order model.OrderResponse
		if err := json.Unmarshal([]byte(cached), &order); err == nil {
			return &order, nil
		}
	}

	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	response := &model.OrderResponse{
		Order: *order,
	}

	// 設置快取
	if cached, err := json.Marshal(response); err == nil {
		s.redis.Set(ctx, cacheKey, cached, 1*time.Hour)
	}

	return response, nil
}

// UpdateOrder 更新訂單
func (s *orderService) UpdateOrder(ctx context.Context, id string, req *model.UpdateOrderRequest) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return ErrOrderNotFound
	}

	if req.Status != "" {
		order.Status = req.Status
	}
	if req.Address != (model.Address{}) {
		order.Address = req.Address
	}
	order.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// 清除快取
	s.clearOrderCache(ctx, order.UserID)
	s.redis.Del(ctx, fmt.Sprintf("order:%s", id))

	return nil
}

// DeleteOrder 刪除訂單
func (s *orderService) DeleteOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return ErrOrderNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	// 清除快取
	s.clearOrderCache(ctx, order.UserID)
	s.redis.Del(ctx, fmt.Sprintf("order:%s", id))

	return nil
}

// ListOrders 獲取訂單列表
func (s *orderService) ListOrders(ctx context.Context, page, limit int) (*model.OrderListResponse, error) {
	orders, total, err := s.repo.List(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	response := &model.OrderListResponse{
		Orders: make([]model.OrderResponse, len(orders)),
		Total:  total,
		Page:   page,
		Limit:  limit,
	}

	for i, order := range orders {
		response.Orders[i] = model.OrderResponse{Order: order}
	}

	return response, nil
}

// GetUserOrders 獲取用戶訂單
func (s *orderService) GetUserOrders(ctx context.Context, userID string, page, limit int) ([]*model.Order, error) {
	orders, _, err := s.repo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	// 轉換 []model.Order 為 []*model.Order
	result := make([]*model.Order, len(orders))
	for i := range orders {
		result[i] = &orders[i]
	}

	return result, nil
}

// GetOrdersByStatus 根據狀態獲取訂單
func (s *orderService) GetOrdersByStatus(ctx context.Context, status model.OrderStatus, page, limit int) (*model.OrderListResponse, error) {
	orders, total, err := s.repo.GetByStatus(ctx, status, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by status: %w", err)
	}

	response := &model.OrderListResponse{
		Orders: make([]model.OrderResponse, len(orders)),
		Total:  total,
		Page:   page,
		Limit:  limit,
	}

	for i, order := range orders {
		response.Orders[i] = model.OrderResponse{Order: order}
	}

	return response, nil
}

// CancelOrder 取消訂單
func (s *orderService) CancelOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return ErrOrderNotFound
	}

	if order.Status != model.OrderStatusPending {
		return ErrInvalidOrderState
	}

	order.Status = model.OrderStatusCancelled
	order.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// 清除快取
	s.clearOrderCache(ctx, order.UserID)
	s.redis.Del(ctx, fmt.Sprintf("order:%s", id))

	return nil
}

// clearOrderCache 清除訂單相關快取
func (s *orderService) clearOrderCache(ctx context.Context, userID string) {
	s.redis.Del(ctx, fmt.Sprintf("user:%s:orders", userID))
}

func (s *orderService) CreateOrderFromCart(ctx context.Context, userID string, req *model.CreateOrderFromCartRequest) (*model.CreateOrderResponse, error) {
	// 計算訂單總金額
	var totalAmount float64
	for _, item := range req.CartItems {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// 創建訂單
	order := &model.Order{
		ID:          uuid.New().String(),
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      model.OrderStatusPending,
		Items:       make([]model.OrderItem, len(req.CartItems)),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 轉換購物車項目為訂單項目

	for i, cartItem := range req.CartItems {
		order.Items[i] = model.OrderItem{
			OrderID:   order.ID,
			ProductID: cartItem.ProductID,
			Price:     cartItem.Price,
			Quantity:  cartItem.Quantity,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	// 保存訂單
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return &model.CreateOrderResponse{
		OrderID:     order.ID,
		TotalAmount: order.TotalAmount,
		Status:      string(order.Status),
		CreatedAt:   order.CreatedAt,
	}, nil
}

// GetOrders 獲取訂單列表
func (s *orderService) GetOrders(ctx context.Context, page, limit int) ([]*model.Order, error) {
	// 嘗試從快取獲取
	cacheKey := fmt.Sprintf("orders:page:%d:limit:%d", page, limit)
	if cached, err := s.redis.Get(ctx, cacheKey); err == nil {
		var orders []*model.Order
		if err := json.Unmarshal([]byte(cached), &orders); err == nil {
			return orders, nil
		}
	}

	// 從資料庫獲取
	orders, err := s.repo.GetOrders(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	// 設置快取
	if cached, err := json.Marshal(orders); err == nil {
		s.redis.Set(ctx, cacheKey, cached, 1*time.Hour)
	}

	return orders, nil
}
