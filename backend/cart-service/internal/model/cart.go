package model

import (
	"time"

	"gorm.io/gorm"
)

// Cart 購物車模型
type Cart struct {
	UserID    string     `gorm:"primaryKey"`
	Items     []CartItem `gorm:"foreignKey:UserID;references:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// CartItem 購物車項目模型
type CartItem struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     string // 直接使用 UserID 來替代 CartUserID
	ProductID  string
	Name       string  // 商品名稱
	Image      string  // 商品圖片 (base64 格式)
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
	ProductID string `json:"ProductID" binding:"required"`
	Quantity  int    `json:"Quantity" binding:"required,gte=0"`
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

// 新增一個用於處理圖片的類型
type ProductImage struct {
	URL  string `json:"url"`  // 原始URL
	Data string `json:"data"` // base64編碼的圖片數據
}
