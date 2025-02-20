package model

import (
	"time"
)

// OrderStatus 訂單狀態
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipping  OrderStatus = "shipping"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order 訂單模型
type Order struct {
	ID          string      `json:"id" gorm:"primaryKey;type:string"`
	UserID      string      `json:"userId" gorm:"index;type:string"`
	Status      OrderStatus `json:"status"`
	TotalAmount float64     `json:"totalAmount"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Address     Address     `json:"address" gorm:"embedded"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	DeletedAt   *time.Time  `json:"deletedAt,omitempty" gorm:"index"`
}

// OrderItem 訂單項目
type OrderItem struct {
	ID        string    `json:"id" gorm:"primaryKey;type:string"`
	OrderID   string    `json:"orderId" gorm:"index;type:string"`
	ProductID string    `json:"productId"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Address 地址
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
}

// CreateOrderRequest 創建訂單請求
type CreateOrderRequest struct {
	UserID  string      `json:"-"`  // 從 JWT token 中獲取，不從 JSON 讀取
	Items   []OrderItem `json:"items" binding:"required,min=1"`
	Address Address     `json:"address" binding:"required"`
}

// UpdateOrderRequest 更新訂單請求
type UpdateOrderRequest struct {
	Status  OrderStatus `json:"status,omitempty"`
	Address Address     `json:"address,omitempty"`
}

// OrderResponse 訂單響應
type OrderResponse struct {
	Order
	PaymentStatus string `json:"paymentStatus,omitempty"`
}

// OrderListResponse 訂單列表響應
type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int64          `json:"total"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
}
