package repository

import (
	"context"
	"fmt"

	"firebase.google.com/go/db"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/model"
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
	db *db.Client
}

// NewPaymentRepository 創建支付存儲實例
func NewPaymentRepository(db *db.Client) PaymentRepository {
	return &paymentRepository{db: db}
}

// Create 創建支付記錄
func (r *paymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	ref := r.db.NewRef("payments")
	return ref.Child(payment.ID).Set(ctx, payment)
}

// GetByID 根據ID獲取支付記錄
func (r *paymentRepository) GetByID(ctx context.Context, id string) (*model.Payment, error) {
	var payment model.Payment
	ref := r.db.NewRef("payments").Child(id)
	if err := ref.Get(ctx, &payment); err != nil {
		return nil, err
	}
	if payment.ID == "" {
		return nil, nil
	}
	return &payment, nil
}

// Update 更新支付記錄
func (r *paymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	ref := r.db.NewRef("payments").Child(payment.ID)
	return ref.Set(ctx, payment)
}

// GetByOrderID 根據訂單ID獲取支付記錄
func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	var result map[string]model.Payment
	ref := r.db.NewRef("payments")
	if err := ref.OrderByChild("orderId").EqualTo(orderID).Get(ctx, &result); err != nil {
		return nil, err
	}
	for _, payment := range result {
		return &payment, nil
	}
	return nil, nil
}

// GetByUserID 獲取用戶的支付記錄
func (r *paymentRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]model.Payment, int64, error) {
	var result map[string]model.Payment
	ref := r.db.NewRef("payments")

	// 獲取所有符合用戶ID的記錄
	if err := ref.OrderByChild("user_id").EqualTo(userID).Get(ctx, &result); err != nil {
		return nil, 0, err
	}

	// 計算總數
	total := int64(len(result))

	// 計算分頁
	start := (page - 1) * limit
	end := start + limit
	if start >= int(total) {
		return []model.Payment{}, total, nil
	}
	if end > int(total) {
		end = int(total)
	}

	// 轉換為切片並進行分頁
	payments := make([]model.Payment, 0, len(result))
	for _, payment := range result {
		payments = append(payments, payment)
	}

	// 返回分頁後的結果
	return payments[start:end], total, nil
}

// List 獲取支付記錄列表
func (r *paymentRepository) List(ctx context.Context, page, limit int) ([]model.Payment, int64, error) {
	var result map[string]model.Payment
	ref := r.db.NewRef("payments")

	if err := ref.Get(ctx, &result); err != nil {
		return nil, 0, err
	}

	total := int64(len(result))

	// 計算分頁
	start := (page - 1) * limit
	end := start + limit
	if start >= int(total) {
		return []model.Payment{}, total, nil
	}
	if end > int(total) {
		end = int(total)
	}

	// 轉換為切片並進行分頁
	payments := make([]model.Payment, 0, len(result))
	for _, payment := range result {
		payments = append(payments, payment)
	}

	return payments[start:end], total, nil
}

// CreateRefund 創建退款記錄
func (r *paymentRepository) CreateRefund(ctx context.Context, refund *model.Refund) error {
	ref := r.db.NewRef(fmt.Sprintf("payments/%s/refunds", refund.PaymentID))
	return ref.Child(refund.ID).Set(ctx, refund)
}

// GetRefundsByPaymentID 獲取支付的退款記錄
func (r *paymentRepository) GetRefundsByPaymentID(ctx context.Context, paymentID string) ([]model.Refund, error) {
	var result map[string]model.Refund
	ref := r.db.NewRef(fmt.Sprintf("payments/%s/refunds", paymentID))

	if err := ref.Get(ctx, &result); err != nil {
		return nil, err
	}

	refunds := make([]model.Refund, 0, len(result))
	for _, refund := range result {
		refunds = append(refunds, refund)
	}

	return refunds, nil
}
