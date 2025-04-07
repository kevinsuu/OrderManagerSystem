package model

import "time"

// OrderStatus 訂單狀態
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order 訂單模型
type Order struct {
	ID           string       `json:"id"`
	UserID       string       `json:"userId"`
	Items        []OrderItem  `json:"items"`
	TotalAmount  float64      `json:"totalAmount"`
	Status       OrderStatus  `json:"status"`
	ShippingInfo ShippingInfo `json:"shippingInfo"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

// OrderItem 訂單項目
type OrderItem struct {
	ProductID  string  `json:"productId"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"totalPrice"`
}

// ShippingInfo 配送信息
type ShippingInfo struct {
	RecipientName  string  `json:"recipientName"`
	PhoneNumber    string  `json:"phoneNumber"`
	Address        Address `json:"address"`
	ShippingMethod string  `json:"shippingMethod"`
}

// Address 地址信息
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}
