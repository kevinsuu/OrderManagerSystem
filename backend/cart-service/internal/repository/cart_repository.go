package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/db"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
)

type CartRepository interface {
	GetCart(ctx context.Context, userID string) (*model.Cart, error)
	SaveCart(ctx context.Context, cart *model.Cart) error
	DeleteCart(ctx context.Context, userID string) error
	UpdateCartItems(ctx context.Context, userID string, items []model.CartItem) error
	AddItem(ctx context.Context, userID string, item model.CartItem) error
	RemoveItem(ctx context.Context, userID string, productID string) error
	UpdateQuantity(ctx context.Context, userID string, productID string, quantity int) error
	SelectItems(ctx context.Context, userID string, productIDs []string) error
	ClearCart(ctx context.Context, userID string) error
}

type cartRepository struct {
	client *db.Client
}

func NewCartRepository(client *db.Client) CartRepository {
	return &cartRepository{
		client: client,
	}
}

func (r *cartRepository) GetCart(ctx context.Context, userID string) (*model.Cart, error) {
	var cart model.Cart
	if err := r.client.NewRef("carts").Child(userID).Get(ctx, &cart); err != nil {
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}
	fmt.Printf("Cart retrieved from repository: %+v", cart)
	return &cart, nil
}

func (r *cartRepository) SaveCart(ctx context.Context, cart *model.Cart) error {
	return r.client.NewRef("carts").Child(cart.UserID).Set(ctx, cart)
}

func (r *cartRepository) DeleteCart(ctx context.Context, userID string) error {
	return r.client.NewRef("carts").Child(userID).Delete(ctx)
}

func (r *cartRepository) UpdateCartItems(ctx context.Context, userID string, items []model.CartItem) error {
	return r.client.NewRef("carts").Child(userID).Child("items").Set(ctx, items)
}

func (r *cartRepository) AddItem(ctx context.Context, userID string, item model.CartItem) error {
	log.Printf("Adding item to cart for user: %s, product: %s", userID, item.ProductID)

	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		log.Printf("Error getting cart, creating new one: %v", err)
		cart = &model.Cart{
			UserID:    userID,
			Items:     []model.CartItem{item},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		log.Printf("New cart created with userID: %s", userID)
		return r.SaveCart(ctx, cart)
	}

	// 確保userID已設置
	if cart.UserID == "" {
		cart.UserID = userID
		log.Printf("Setting userID for existing cart: %s", userID)
	}

	// 更新時間戳
	cart.UpdatedAt = time.Now()

	// 檢查是否已有該商品
	found := false
	for i, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			cart.Items[i].Quantity += item.Quantity
			cart.Items[i].UpdatedAt = time.Now()
			found = true
			log.Printf("Updated existing item quantity to: %d", cart.Items[i].Quantity)
			break
		}
	}

	// 如果沒有找到該商品，添加新項目
	if !found {
		cart.Items = append(cart.Items, item)
		log.Printf("Added new item to cart, total items: %d", len(cart.Items))
	}

	// 保存購物車
	err = r.SaveCart(ctx, cart)
	if err != nil {
		log.Printf("Error saving cart: %v", err)
		return err
	}

	log.Printf("Cart saved successfully for user: %s", userID)
	return nil
}

func (r *cartRepository) RemoveItem(ctx context.Context, userID string, productID string) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	var updatedItems []model.CartItem
	for _, item := range cart.Items {
		if item.ProductID != productID {
			updatedItems = append(updatedItems, item)
		}
	}

	cart.Items = updatedItems
	return r.SaveCart(ctx, cart)
}

func (r *cartRepository) UpdateQuantity(ctx context.Context, userID string, productID string, quantity int) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].Quantity = quantity
			return r.SaveCart(ctx, cart)
		}
	}

	return fmt.Errorf("product not found in cart")
}

func (r *cartRepository) SelectItems(ctx context.Context, userID string, productIDs []string) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	selectedProducts := make(map[string]bool)
	for _, id := range productIDs {
		selectedProducts[id] = true
	}

	for i := range cart.Items {
		cart.Items[i].Selected = selectedProducts[cart.Items[i].ProductID]
	}

	return r.SaveCart(ctx, cart)
}

func (r *cartRepository) ClearCart(ctx context.Context, userID string) error {
	return r.DeleteCart(ctx, userID)
}
