package model

import (
	"time"
)

// Cart 購物車模型
type Cart struct {
	UserID    string     `json:"userId"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// CartItem 購物車項目
type CartItem struct {
	ProductID  string    `json:"productId"`
	Name       string    `json:"name"`       // 商品名稱
	Image      string    `json:"image"`      // 商品圖片
	Price      float64   `json:"price"`      // 當前價格
	Quantity   int       `json:"quantity"`   // 數量
	Selected   bool      `json:"selected"`   // 是否選中
	StockCount int       `json:"stockCount"` // 當前庫存數量
	UpdatedAt  time.Time `json:"updatedAt"`
}

// AddToCartRequest 添加商品到購物車請求
type AddToCartRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

// UpdateQuantityRequest 更新商品數量請求
type UpdateQuantityRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gte=0"`
}

// SelectItemsRequest 選擇商品請求
type SelectItemsRequest struct {
	ProductIDs []string `json:"productIds" binding:"required"`
}

// CartResponse 購物車響應
type CartResponse struct {
	Items         []CartItem `json:"items"`
	TotalSelected int        `json:"totalSelected"` // 已選商品總數
	TotalAmount   float64    `json:"totalAmount"`   // 已選商品總金額
}
