package repository

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/db"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
)

// WishlistRepository 提供收藏清單相關操作
type WishlistRepository interface {
	AddToWishlist(ctx context.Context, userId, productId string) error
	RemoveFromWishlist(ctx context.Context, userId, productId string) error
	GetWishlist(ctx context.Context, userId string, page, limit int) (*model.WishlistResponse, error)
	IsProductInWishlist(ctx context.Context, userId, productId string) (bool, error)
}

// wishlistRepository 實現 WishlistRepository 接口
type wishlistRepository struct {
	db *db.Client
}

// NewWishlistRepository 創建一個新的收藏清單儲存庫
func NewWishlistRepository(db *db.Client) WishlistRepository {
	return &wishlistRepository{db: db}
}

// AddToWishlist 添加商品到收藏清單
func (r *wishlistRepository) AddToWishlist(ctx context.Context, userId, productId string) error {
	if userId == "" || productId == "" {
		return fmt.Errorf("userId and productId cannot be empty")
	}

	// 建立收藏項目
	wishlistItem := model.WishlistItem{
		UserId:    userId,
		ProductId: productId,
		CreatedAt: time.Now(),
		ID:        userId + "_" + productId,
	}

	// 儲存到 Firebase Realtime Database
	key := fmt.Sprintf("%s_%s", userId, productId)
	ref := r.db.NewRef("wishlists").Child(key)
	return ref.Set(ctx, wishlistItem)
}

// RemoveFromWishlist 從收藏清單移除商品
func (r *wishlistRepository) RemoveFromWishlist(ctx context.Context, userId, productId string) error {
	if userId == "" || productId == "" {
		return fmt.Errorf("userId and productId cannot be empty")
	}

	// 從 Firebase Realtime Database 刪除
	key := fmt.Sprintf("%s_%s", userId, productId)
	ref := r.db.NewRef("wishlists").Child(key)
	return ref.Delete(ctx)
}

// GetWishlist 獲取使用者的收藏清單
func (r *wishlistRepository) GetWishlist(ctx context.Context, userId string, page, limit int) (*model.WishlistResponse, error) {
	if userId == "" {
		return nil, fmt.Errorf("userId cannot be empty")
	}

	// 計算分頁
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// 查詢使用者的收藏項目
	ref := r.db.NewRef("wishlists")
	// Firebase RTDB 沒有直接的 WHERE 查詢，所以我們需要獲取所有數據然後在內存中過濾
	var allItems map[string]model.WishlistItem
	if err := ref.Get(ctx, &allItems); err != nil {
		// 如果是數據不存在的錯誤，返回空結果
		if err.Error() == "client: response error: data at path \"wishlists\" does not exist" {
			return &model.WishlistResponse{
				Wishlist: []model.WishlistItem{},
				Total:    0,
				Page:     page,
				Limit:    limit,
			}, nil
		}
		return nil, fmt.Errorf("failed to get wishlist items: %w", err)
	}

	// 過濾出特定用戶的收藏
	var userItems []model.WishlistItem
	for key, item := range allItems {
		if item.UserId == userId {
			// 使用Firebase生成的完整key作為ID
			item.ID = key
			userItems = append(userItems, item)
		}
	}

	// 計算總數
	total := len(userItems)

	// 排序 - 按創建時間降序（最新的排在前面）
	if len(userItems) > 1 {
		// 使用冒泡排序進行簡單排序
		for i := 0; i < len(userItems)-1; i++ {
			for j := 0; j < len(userItems)-i-1; j++ {
				if userItems[j].CreatedAt.Before(userItems[j+1].CreatedAt) {
					userItems[j], userItems[j+1] = userItems[j+1], userItems[j]
				}
			}
		}
	}

	// 分頁
	start := (page - 1) * limit
	end := start + limit
	if start >= total {
		start = 0
		end = 0
	}
	if end > total {
		end = total
	}

	var pagedItems []model.WishlistItem
	if start < end {
		pagedItems = userItems[start:end]
	}

	// 返回收藏清單和總數
	return &model.WishlistResponse{
		Wishlist: pagedItems,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

// IsProductInWishlist 檢查商品是否在收藏清單中
func (r *wishlistRepository) IsProductInWishlist(ctx context.Context, userId, productId string) (bool, error) {
	if userId == "" || productId == "" {
		return false, fmt.Errorf("userId and productId cannot be empty")
	}

	// 檢查商品是否在收藏清單中
	ref := r.db.NewRef("wishlists")
	key := fmt.Sprintf("%s_%s", userId, productId)
	var item model.WishlistItem
	if err := ref.Child(key).Get(ctx, &item); err != nil {
		// 如果是數據不存在的錯誤，返回false表示不在清單中
		if err.Error() == "client: response error: data at path \"wishlists/"+key+"\" does not exist" {
			return false, nil
		}
		return false, fmt.Errorf("failed to get wishlist item: %w", err)
	}

	// 如果找到了項目且ProductId不為空，則表示商品在收藏清單中
	return item.ProductId != "", nil
}
