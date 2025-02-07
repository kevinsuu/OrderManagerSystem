package repository

import (
	"context"
	"errors"

	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/model"
	"gorm.io/gorm"
)

// PaymentRepository 支付存儲接口
type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
	GetByID(ctx context.Context, id string) (*model.Payment, error)
	Update(ctx context.Context, payment *model.Payment) error
	GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error)
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]model.Payment, int64, error)
	List(ctx context.Context, page, limit int) ([]model.Payment, int64, error)
	CreateRefund(ctx context.Context, refund *model.Refund) error
	GetRefundsByPaymentID(ctx context.Context, paymentID string) ([]model.Refund, error)
}

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository 創建支付存儲實例
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// Create 創建支付記錄
func (r *paymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// GetByID 根據ID獲取支付記錄
func (r *paymentRepository) GetByID(ctx context.Context, id string) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.WithContext(ctx).First(&payment, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// Update 更新支付記錄
func (r *paymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

// GetByOrderID 根據訂單ID獲取支付記錄
func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.WithContext(ctx).First(&payment, "order_id = ?", orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// GetByUserID 獲取用戶的支付記錄
func (r *paymentRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]model.Payment, int64, error) {
	var payments []model.Payment
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&model.Payment{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// List 獲取支付記錄列表
func (r *paymentRepository) List(ctx context.Context, page, limit int) ([]model.Payment, int64, error) {
	var payments []model.Payment
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&model.Payment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// CreateRefund 創建退款記錄
func (r *paymentRepository) CreateRefund(ctx context.Context, refund *model.Refund) error {
	return r.db.WithContext(ctx).Create(refund).Error
}

// GetRefundsByPaymentID 獲取支付的退款記錄
func (r *paymentRepository) GetRefundsByPaymentID(ctx context.Context, paymentID string) ([]model.Refund, error) {
	var refunds []model.Refund
	if err := r.db.WithContext(ctx).
		Where("payment_id = ?", paymentID).
		Find(&refunds).Error; err != nil {
		return nil, err
	}
	return refunds, nil
}
