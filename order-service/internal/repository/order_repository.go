package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kevinsuu/OrderManagerSystem/order-service/internal/model"
	"gorm.io/gorm"
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
	db *gorm.DB
}

// NewOrderRepository 創建訂單存儲實例
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create 創建訂單
func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	// 生成訂單 ID
	order.ID = uuid.New().String()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// 使用事務來確保訂單和訂單項目的一致性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先創建訂單（不包含 Items）
		if err := tx.Omit("Items").Create(order).Error; err != nil {
			return err
		}

		// 準備訂單項目數據
		for i := range order.Items {
			order.Items[i].ID = uuid.New().String()
			order.Items[i].OrderID = order.ID
			order.Items[i].CreatedAt = time.Now()
			order.Items[i].UpdatedAt = time.Now()
		}

		// 分批創建訂單項目
		if len(order.Items) > 0 {
			// 使用 CreateInBatches 來批量插入
			if err := tx.CreateInBatches(order.Items, 100).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetByID 根據ID獲取訂單
func (r *orderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	if err := r.db.WithContext(ctx).Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// Update 更新訂單
func (r *orderRepository) Update(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete 刪除訂單
func (r *orderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Order{}, "id = ?", id).Error
}

// List 獲取訂單列表
func (r *orderRepository) List(ctx context.Context, page, limit int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&model.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Items").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetByUserID 獲取用戶的訂單
func (r *orderRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetByStatus 根據狀態獲取訂單
func (r *orderRepository) GetByStatus(ctx context.Context, status model.OrderStatus, page, limit int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&model.Order{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Items").Where("status = ?", status).Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetOrder 根據ID獲取訂單
func (r *orderRepository) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	if err := r.db.WithContext(ctx).Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// GetOrders 獲取訂單列表
func (r *orderRepository) GetOrders(ctx context.Context, page, limit int) ([]*model.Order, error) {
	var orders []*model.Order
	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).
		Preload("Items").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// GetUserOrders 獲取用戶的訂單列表
func (r *orderRepository) GetUserOrders(ctx context.Context, userID string, page, limit int) ([]*model.Order, error) {
	var orders []*model.Order
	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// UpdateOrder 更新訂單
func (r *orderRepository) UpdateOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}
