package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/db"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
)

// OrderRepository 訂單倉庫接口
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, orderID string) (*model.Order, error)
	GetByUserID(ctx context.Context, userID string, offset, limit int) ([]model.Order, error)
	UpdateStatus(ctx context.Context, orderID string, status model.OrderStatus) error
}

type orderRepository struct {
	client *db.Client
}

// NewOrderRepository 創建訂單倉庫實例
func NewOrderRepository(client *db.Client) OrderRepository {
	return &orderRepository{
		client: client,
	}
}

// Create 創建訂單
func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	// 設置創建和更新時間
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	// 保存訂單
	return r.client.NewRef("orders").Child(order.ID).Set(ctx, order)
}

// GetByID 根據ID獲取訂單
func (r *orderRepository) GetByID(ctx context.Context, orderID string) (*model.Order, error) {
	var order model.Order
	if err := r.client.NewRef("orders").Child(orderID).Get(ctx, &order); err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}
	return &order, nil
}

// GetByUserID 根據用戶ID獲取訂單
func (r *orderRepository) GetByUserID(ctx context.Context, userID string, offset, limit int) ([]model.Order, error) {
	var orders map[string]model.Order
	if err := r.client.NewRef("orders").OrderByChild("userId").EqualTo(userID).Get(ctx, &orders); err != nil {
		log.Printf("Error getting orders by user ID: %v", err)
		return nil, fmt.Errorf("error getting orders by user ID: %v", err)
	}

	result := make([]model.Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, order)
	}

	// 應用分頁
	start := offset
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	if start >= len(result) {
		return []model.Order{}, nil
	}

	return result[start:end], nil
}

// UpdateStatus 更新訂單狀態
func (r *orderRepository) UpdateStatus(ctx context.Context, orderID string, status model.OrderStatus) error {
	updates := map[string]interface{}{
		"status":    status,
		"updatedAt": time.Now(),
	}
	return r.client.NewRef("orders").Child(orderID).Update(ctx, updates)
}
