package model

import (
	"time"
)

// PaymentStatus 支付狀態
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PaymentMethod 支付方式
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodDigitalWallet PaymentMethod = "digital_wallet"
)

// Payment 支付模型
type Payment struct {
	ID            string        `json:"id" gorm:"primaryKey"`
	OrderID       string        `json:"orderId" gorm:"index"`
	UserID        string        `json:"userId" gorm:"index"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	Status        PaymentStatus `json:"status"`
	Method        PaymentMethod `json:"method"`
	TransactionID string        `json:"transactionId"`
	ErrorMessage  string        `json:"errorMessage,omitempty"`
	Metadata      string        `json:"metadata,omitempty"` // JSON 字符串，存儲額外信息
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	DeletedAt     *time.Time    `json:"deletedAt,omitempty" gorm:"index"`
}

// Refund 退款模型
type Refund struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	PaymentID     string    `json:"paymentId" gorm:"index"`
	Amount        float64   `json:"amount"`
	Reason        string    `json:"reason"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transactionId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// CreatePaymentRequest 創建支付請求
type CreatePaymentRequest struct {
	OrderID  string        `json:"orderId" binding:"required"`
	UserID   string        `json:"userId" binding:"required"`
	Amount   float64       `json:"amount" binding:"required,gt=0"`
	Currency string        `json:"currency" binding:"required"`
	Method   PaymentMethod `json:"method" binding:"required"`
}

// RefundRequest 退款請求
type RefundRequest struct {
	PaymentID string  `json:"paymentId" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Reason    string  `json:"reason" binding:"required"`
}

// PaymentResponse 支付響應
type PaymentResponse struct {
	Payment
	RefundHistory []Refund `json:"refundHistory,omitempty"`
}

// PaymentListResponse 支付列表響應
type PaymentListResponse struct {
	Payments []PaymentResponse `json:"payments"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

// PaymentGatewayResponse 支付網關響應
type PaymentGatewayResponse struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transactionId,omitempty"`
	ErrorMessage  string `json:"errorMessage,omitempty"`
}

// PaymentNotification 支付通知
type PaymentNotification struct {
	PaymentID     string        `json:"paymentId"`
	OrderID       string        `json:"orderId"`
	Status        PaymentStatus `json:"status"`
	TransactionID string        `json:"transactionId"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	Timestamp     time.Time     `json:"timestamp"`
}
