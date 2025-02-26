package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetCart(ctx context.Context, userID string) (*model.Cart, error)
	AddItem(ctx context.Context, userID string, item model.CartItem) error
	RemoveItem(ctx context.Context, userID string, productID string) error
	UpdateQuantity(ctx context.Context, userID string, productID string, quantity int) error
	ClearCart(ctx context.Context, userID string) error
	SelectItems(ctx context.Context, userID string, productIDs []string) error
}

type cartRepository struct {
	redis RedisRepository
	db    *gorm.DB
}

func NewCartRepository(redis RedisRepository, db *gorm.DB) CartRepository {
	return &cartRepository{
		redis: redis,
		db:    db,
	}
}

func (r *cartRepository) GetCart(ctx context.Context, userID string) (*model.Cart, error) {
	key := fmt.Sprintf("cart:%s", userID)
	data, err := r.redis.Get(ctx, key)
	if err != nil {
		return &model.Cart{
			UserID: userID,
			Items:  []model.CartItem{},
		}, nil
	}

	var cart model.Cart
	if err := json.Unmarshal([]byte(data), &cart); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) AddItem(ctx context.Context, userID string, item model.CartItem) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// 檢查商品是否已存在
	found := false
	for i, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			cart.Items[i].Quantity += item.Quantity
			cart.Items[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		item.UpdatedAt = time.Now()
		cart.Items = append(cart.Items, item)
	}

	cart.UpdatedAt = time.Now()

	// 保存更新後的購物車
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("cart:%s", userID)
	return r.redis.Set(ctx, key, data, 24*time.Hour) // 購物車數據24小時過期
}

func (r *cartRepository) RemoveItem(ctx context.Context, userID string, productID string) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// 找到並移除商品
	newItems := make([]model.CartItem, 0)
	for _, item := range cart.Items {
		if item.ProductID != productID {
			newItems = append(newItems, item)
		}
	}
	cart.Items = newItems
	cart.UpdatedAt = time.Now()

	// 保存更新後的購物車
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("cart:%s", userID)
	return r.redis.Set(ctx, key, data, 24*time.Hour)
}

func (r *cartRepository) UpdateQuantity(ctx context.Context, userID string, productID string, quantity int) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// 更新商品數量
	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].Quantity = quantity
			cart.Items[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("product not found in cart")
	}

	cart.UpdatedAt = time.Now()

	// 保存更新後的購物車
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("cart:%s", userID)
	return r.redis.Set(ctx, key, data, 24*time.Hour)
}

func (r *cartRepository) SelectItems(ctx context.Context, userID string, productIDs []string) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// 創建一個 map 來存儲要選擇的商品 ID
	selectedProducts := make(map[string]bool)
	for _, id := range productIDs {
		selectedProducts[id] = true
	}

	// 更新商品的選中狀態
	for i := range cart.Items {
		cart.Items[i].Selected = selectedProducts[cart.Items[i].ProductID]
		if cart.Items[i].Selected {
			cart.Items[i].UpdatedAt = time.Now()
		}
	}

	cart.UpdatedAt = time.Now()

	// 保存更新後的購物車
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("cart:%s", userID)
	return r.redis.Set(ctx, key, data, 24*time.Hour)
}

func (r *cartRepository) ClearCart(ctx context.Context, userID string) error {
	key := fmt.Sprintf("cart:%s", userID)
	return r.redis.Del(ctx, key)
}
