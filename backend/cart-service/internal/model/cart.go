package model

import (
	"time"

	"gorm.io/gorm"
)

// Cart 購物車模型
type Cart struct {
	UserID    string     `gorm:"primaryKey"`
	Items     []CartItem `gorm:"foreignKey:CartUserID;references:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// CartItem 購物車項目模型
type CartItem struct {
	ID         uint   `gorm:"primaryKey"`
	CartUserID string // 外鍵，關聯到 Cart 的 UserID
	ProductID  string
	Name       string  // 商品名稱
	Image      string  // 商品圖片
	Price      float64 // 商品價格
	StockCount int     // 庫存數量
	Quantity   int
	Selected   bool `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
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
