package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/repository"
)

var (
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrInvalidPaymentStatus = errors.New("invalid payment status")
	ErrInvalidRefundAmount  = errors.New("invalid refund amount")
)

// PaymentService 支付服務接口
type PaymentService interface {
	CreatePayment(ctx context.Context, req *model.CreatePaymentRequest) (*model.Payment, error)
	GetPayment(ctx context.Context, id string) (*model.PaymentResponse, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*model.PaymentResponse, error)
	GetUserPayments(ctx context.Context, userID string, page, limit int) (*model.PaymentListResponse, error)
	ListPayments(ctx context.Context, page, limit int) (*model.PaymentListResponse, error)
	ProcessPayment(ctx context.Context, id string) error
	RefundPayment(ctx context.Context, req *model.RefundRequest) error
	CancelPayment(ctx context.Context, id string) error
}

type paymentService struct {
	repo  repository.PaymentRepository
	redis repository.RedisRepository
}

// NewPaymentService 創建支付服務實例
func NewPaymentService(repo repository.PaymentRepository, redis repository.RedisRepository) PaymentService {
	return &paymentService{
		repo:  repo,
		redis: redis,
	}
}

// CreatePayment 創建支付
func (s *paymentService) CreatePayment(ctx context.Context, req *model.CreatePaymentRequest) (*model.Payment, error) {
	payment := &model.Payment{
		ID:        uuid.New().String(),
		OrderID:   req.OrderID,
		UserID:    req.UserID,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Status:    model.PaymentStatusPending,
		Method:    req.Method,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// 清除相關快取
	s.clearPaymentCache(ctx, payment.UserID)

	return payment, nil
}

// GetPayment 獲取支付詳情
func (s *paymentService) GetPayment(ctx context.Context, id string) (*model.PaymentResponse, error) {
	// 嘗試從快取獲取
	cacheKey := fmt.Sprintf("payment:%s", id)
	if cached, err := s.redis.Get(ctx, cacheKey); err == nil {
		var response model.PaymentResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	// 獲取退款歷史
	refunds, err := s.repo.GetRefundsByPaymentID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get refund history: %w", err)
	}

	response := &model.PaymentResponse{
		Payment:       *payment,
		RefundHistory: refunds,
	}

	// 設置快取
	if cached, err := json.Marshal(response); err == nil {
		s.redis.Set(ctx, cacheKey, cached, 1*time.Hour)
	}

	return response, nil
}

// GetPaymentByOrderID 根據訂單ID獲取支付
func (s *paymentService) GetPaymentByOrderID(ctx context.Context, orderID string) (*model.PaymentResponse, error) {
	payment, err := s.repo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by order ID: %w", err)
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	refunds, err := s.repo.GetRefundsByPaymentID(ctx, payment.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get refund history: %w", err)
	}

	return &model.PaymentResponse{
		Payment:       *payment,
		RefundHistory: refunds,
	}, nil
}

// GetUserPayments 獲取用戶支付記錄
func (s *paymentService) GetUserPayments(ctx context.Context, userID string, page, limit int) (*model.PaymentListResponse, error) {
	payments, total, err := s.repo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user payments: %w", err)
	}

	response := &model.PaymentListResponse{
		Payments: make([]model.PaymentResponse, len(payments)),
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	for i, payment := range payments {
		refunds, err := s.repo.GetRefundsByPaymentID(ctx, payment.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get refund history: %w", err)
		}
		response.Payments[i] = model.PaymentResponse{
			Payment:       payment,
			RefundHistory: refunds,
		}
	}

	return response, nil
}

// ListPayments 獲取支付列表
func (s *paymentService) ListPayments(ctx context.Context, page, limit int) (*model.PaymentListResponse, error) {
	payments, total, err := s.repo.List(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	response := &model.PaymentListResponse{
		Payments: make([]model.PaymentResponse, len(payments)),
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	for i, payment := range payments {
		refunds, err := s.repo.GetRefundsByPaymentID(ctx, payment.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get refund history: %w", err)
		}
		response.Payments[i] = model.PaymentResponse{
			Payment:       payment,
			RefundHistory: refunds,
		}
	}

	return response, nil
}

// ProcessPayment 處理支付
func (s *paymentService) ProcessPayment(ctx context.Context, id string) error {
	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return ErrPaymentNotFound
	}

	if payment.Status != model.PaymentStatusPending {
		return ErrInvalidPaymentStatus
	}

	// 模擬調用支付網關
	gatewayResponse := s.processPaymentGateway(payment)

	payment.TransactionID = gatewayResponse.TransactionID
	payment.Status = model.PaymentStatusSuccess
	if !gatewayResponse.Success {
		payment.Status = model.PaymentStatusFailed
		payment.ErrorMessage = gatewayResponse.ErrorMessage
	}
	payment.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// 清除快取
	s.clearPaymentCache(ctx, payment.UserID)
	s.redis.Del(ctx, fmt.Sprintf("payment:%s", id))

	return nil
}

// RefundPayment 退款
func (s *paymentService) RefundPayment(ctx context.Context, req *model.RefundRequest) error {
	payment, err := s.repo.GetByID(ctx, req.PaymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return ErrPaymentNotFound
	}

	if payment.Status != model.PaymentStatusSuccess {
		return ErrInvalidPaymentStatus
	}

	if req.Amount > payment.Amount {
		return ErrInvalidRefundAmount
	}

	// 創建退款記錄
	refund := &model.Refund{
		ID:        uuid.New().String(),
		PaymentID: payment.ID,
		Amount:    req.Amount,
		Reason:    req.Reason,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 模擬退款處理
	refund.Status = "success"
	refund.TransactionID = uuid.New().String()

	if err := s.repo.CreateRefund(ctx, refund); err != nil {
		return fmt.Errorf("failed to create refund: %w", err)
	}

	// 更新支付狀態
	if req.Amount == payment.Amount {
		payment.Status = model.PaymentStatusRefunded
		if err := s.repo.Update(ctx, payment); err != nil {
			return fmt.Errorf("failed to update payment: %w", err)
		}
	}

	// 清除快取
	s.clearPaymentCache(ctx, payment.UserID)
	s.redis.Del(ctx, fmt.Sprintf("payment:%s", payment.ID))

	return nil
}

// CancelPayment 取消支付
func (s *paymentService) CancelPayment(ctx context.Context, id string) error {
	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return ErrPaymentNotFound
	}

	if payment.Status != model.PaymentStatusPending {
		return ErrInvalidPaymentStatus
	}

	payment.Status = model.PaymentStatusCancelled
	payment.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// 清除快取
	s.clearPaymentCache(ctx, payment.UserID)
	s.redis.Del(ctx, fmt.Sprintf("payment:%s", id))

	return nil
}

// processPaymentGateway 處理支付網關請求（模擬）
func (s *paymentService) processPaymentGateway(payment *model.Payment) *model.PaymentGatewayResponse {
	// 這裡應該實現實際的支付網關邏輯
	// 目前僅作為示例，模擬成功響應
	return &model.PaymentGatewayResponse{
		Success:       true,
		TransactionID: uuid.New().String(),
	}
}

// clearPaymentCache 清除支付相關快取
func (s *paymentService) clearPaymentCache(ctx context.Context, userID string) {
	s.redis.Del(ctx, fmt.Sprintf("user:%s:payments", userID))
}
