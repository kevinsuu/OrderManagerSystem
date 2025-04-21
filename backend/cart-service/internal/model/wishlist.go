package model

import (
	"time"
)

// ProductInfo 表示商品的基本資訊
type ProductInfo struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Price     float64  `json:"price"`
	Images    []string `json:"images,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	UpdatedAt string   `json:"updatedAt,omitempty"`
}

// WishlistItem 表示收藏清單中的商品項目
type WishlistItem struct {
	ID        string       `json:"id" firestore:"-"`                // 文檔ID (不儲存)
	UserId    string       `json:"userId" firestore:"userId"`       // 使用者ID
	ProductId string       `json:"productId" firestore:"productId"` // 商品ID
	CreatedAt time.Time    `json:"createdAt" firestore:"createdAt"` // 添加時間
	Product   *ProductInfo `json:"product,omitempty" firestore:"-"` // 商品詳細資訊 (不儲存)
}

// WishlistResponse 表示獲取收藏清單的回應
type WishlistResponse struct {
	Wishlist []WishlistItem `json:"wishlist"` // 收藏清單項目
	Total    int            `json:"total"`    // 總數
	Page     int            `json:"page"`     // 當前頁碼
	Limit    int            `json:"limit"`    // 每頁數量
}

// WishlistIdsResponse 表示獲取收藏清單ID的回應
type WishlistIdsResponse struct {
	WishlistIds []string `json:"wishlistIds"` // 收藏的商品ID列表
}

// AddToWishlistRequest 添加到收藏清單的請求
type AddToWishlistRequest struct {
	ProductId string `json:"productId" binding:"required"` // 商品ID
}
