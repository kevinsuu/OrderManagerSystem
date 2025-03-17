package repository

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/db"
	"github.com/google/uuid"
	"github.com/kevinsuu/OrderManagerSystem/order-service/internal/model"
)

// OrderRepository 訂單存儲接口
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]model.Order, int64, error)
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]model.Order, int64, error)
	GetByStatus(ctx context.Context, status model.OrderStatus, page, limit int) ([]model.Order, int64, error)
	GetOrders(ctx context.Context, page, limit int) ([]*model.Order, error)
	GetOrder(ctx context.Context, id string) (*model.Order, error)
	GetUserOrders(ctx context.Context, userID string, page, limit int) ([]*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
	CreateOrder(ctx context.Context, order *model.Order) error
}

type orderRepository struct {
	db *db.Client
}

// NewOrderRepository 創建訂單存儲實例
func NewOrderRepository(db *db.Client) OrderRepository {
	return &orderRepository{db: db}
}

// Create 創建訂單
func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	ref := r.db.NewRef("orders").Child(order.ID)
	return ref.Set(ctx, order)
}

// GetByID 根據ID獲取訂單
func (r *orderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	if err := r.db.NewRef("orders").Child(id).Get(ctx, &order); err != nil {
		if err.Error() == "http error status: 404; reason: Permission denied" {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting order: %v", err)
	}

	// 檢查是否為空數據
	if order.ID == "" {
		return nil, nil
	}

	return &order, nil
}

// Update 更新訂單
func (r *orderRepository) Update(ctx context.Context, order *model.Order) error {
	order.UpdatedAt = time.Now()
	return r.db.NewRef("orders").Child(order.ID).Set(ctx, order)
}

// Delete 刪除訂單
func (r *orderRepository) Delete(ctx context.Context, id string) error {
	return r.db.NewRef("orders").Child(id).Delete(ctx)
}

// List 獲取訂單列表
func (r *orderRepository) List(ctx context.Context, page, limit int) ([]model.Order, int64, error) {
	var orders map[string]model.Order
	if err := r.db.NewRef("orders").Get(ctx, &orders); err != nil {
		return nil, 0, fmt.Errorf("error getting orders: %v", err)
	}

	result := make([]model.Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, order)
	}

	total := int64(len(result))
	start := (page - 1) * limit
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	if start >= len(result) {
		return []model.Order{}, total, nil
	}

	return result[start:end], total, nil
}

// GetByUserID 根據用戶ID獲取訂單
func (r *orderRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]model.Order, int64, error) {
	var orders map[string]model.Order
	if err := r.db.NewRef("orders").OrderByChild("userId").EqualTo(userID).Get(ctx, &orders); err != nil {
		return nil, 0, fmt.Errorf("error getting orders by user ID: %v", err)
	}

	result := make([]model.Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, order)
	}

	total := int64(len(result))
	start := (page - 1) * limit
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	if start >= len(result) {
		return []model.Order{}, total, nil
	}

	return result[start:end], total, nil
}

// GetByStatus 根據狀態獲取訂單
func (r *orderRepository) GetByStatus(ctx context.Context, status model.OrderStatus, page, limit int) ([]model.Order, int64, error) {
	var orders map[string]model.Order
	if err := r.db.NewRef("orders").OrderByChild("status").EqualTo(string(status)).Get(ctx, &orders); err != nil {
		return nil, 0, fmt.Errorf("error getting orders by status: %v", err)
	}

	result := make([]model.Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, order)
	}

	total := int64(len(result))
	start := (page - 1) * limit
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	if start >= len(result) {
		return []model.Order{}, total, nil
	}

	return result[start:end], total, nil
}

// GetOrders 獲取所有訂單
func (r *orderRepository) GetOrders(ctx context.Context, page, limit int) ([]*model.Order, error) {
	var orders map[string]*model.Order
	if err := r.db.NewRef("orders").Get(ctx, &orders); err != nil {
		return nil, fmt.Errorf("error getting orders: %v", err)
	}

	result := make([]*model.Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, order)
	}

	start := (page - 1) * limit
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	if start >= len(result) {
		return []*model.Order{}, nil
	}

	return result[start:end], nil
}

// GetOrder 獲取單個訂單
func (r *orderRepository) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	return r.GetByID(ctx, id)
}

// GetUserOrders 獲取用戶訂單
func (r *orderRepository) GetUserOrders(ctx context.Context, userID string, page, limit int) ([]*model.Order, error) {
	var orders map[string]*model.Order
	if err := r.db.NewRef("orders").OrderByChild("userId").EqualTo(userID).Get(ctx, &orders); err != nil {
		return nil, fmt.Errorf("error getting user orders: %v", err)
	}

	result := make([]*model.Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, order)
	}

	start := (page - 1) * limit
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	if start >= len(result) {
		return []*model.Order{}, nil
	}

	return result[start:end], nil
}

// UpdateOrder 更新訂單
func (r *orderRepository) UpdateOrder(ctx context.Context, order *model.Order) error {
	return r.Update(ctx, order)
}

// CreateOrder 創建訂單
func (r *orderRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	// 為訂單項目生成ID
	for i := range order.Items {
		if order.Items[i].ID == "" {
			order.Items[i].ID = uuid.New().String()
		}
	}
	return r.Create(ctx, order)
}
