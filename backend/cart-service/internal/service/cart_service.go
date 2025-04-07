package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/client"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/repository"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidStock    = errors.New("invalid stock quantity")
)

// Add a Config type
type CartServiceConfig struct {
	ProductServiceBaseURL string
}

type CartService interface {
	GetCart(ctx context.Context, userID string) (*model.CartResponse, error)
	GetCartItems(ctx context.Context, userID string) ([]model.CartItem, error)
	AddItem(ctx context.Context, userID string, req *model.AddToCartRequest) error
	RemoveItem(ctx context.Context, userID string, productID string) error
	UpdateQuantity(ctx context.Context, userID string, req *model.UpdateQuantityRequest) error
	ClearCart(ctx context.Context, userID string) error
	SelectItems(ctx context.Context, userID string, req *model.SelectItemsRequest) error
	CreateOrder(ctx context.Context, userID string) error
}

type cartService struct {
	cartRepo      repository.CartRepository
	productClient client.ProductClient
	orderClient   client.OrderClient
	config        *CartServiceConfig // Add config field
}

// Update the constructor
func NewCartService(cartRepo repository.CartRepository, productClient client.ProductClient, orderClient client.OrderClient, config *CartServiceConfig) CartService {
	// If config is nil, provide default values
	if config == nil {
		config = &CartServiceConfig{
			ProductServiceBaseURL: "https://ordermanagersystem-product-service.onrender.com",
		}
	}

	return &cartService{
		cartRepo:      cartRepo,
		productClient: productClient,
		orderClient:   orderClient,
		config:        config,
	}
}

func (s *cartService) GetCart(ctx context.Context, userID string) (*model.CartResponse, error) {

	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := &model.CartResponse{
		Items: cart.Items,
	}

	// 計算已選商品的總數和總金額
	for _, item := range cart.Items {
		if item.Selected {
			response.TotalSelected += item.Quantity
			response.TotalAmount += float64(item.Quantity) * item.Price
		}
	}

	return response, nil
}

// 添加一個輔助函數來截斷字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func (s *cartService) AddItem(ctx context.Context, userID string, req *model.AddToCartRequest) error {
	// 添加日誌
	log.Printf("Attempting to add product %s with quantity %d for user %s", req.ProductID, req.Quantity, userID)

	// 確保userID不為空
	if userID == "" {
		log.Printf("Error: userID is empty")
		return fmt.Errorf("user ID cannot be empty")
	}

	// 調用 product service 檢查商品是否存在及庫存
	productInfo, err := s.productClient.GetProduct(ctx, req.ProductID)
	if err != nil {
		log.Printf("Error getting product info: %v", err)
		return fmt.Errorf("failed to get product info: %w", err)
	}
	if productInfo == nil {
		log.Printf("Product not found: %s", req.ProductID)
		return ErrProductNotFound
	}

	// 添加庫存檢查的日誌
	log.Printf("Product info received: %+v", productInfo)
	log.Printf("Product %s stock count: %d, requested quantity: %d", req.ProductID, productInfo.Stock, req.Quantity)

	// 檢查庫存
	if productInfo.Stock < req.Quantity {
		return fmt.Errorf("insufficient stock: available %d, requested %d", productInfo.Stock, req.Quantity)
	}

	// 檢查購物車中是否已有該商品
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		log.Printf("Error getting cart: %v", err)
		// 這裡不直接返回錯誤，因為GetCart已經處理了創建新購物車的邏輯
	}

	// 計算購物車中已有的數量
	existingQuantity := 0
	for _, item := range cart.Items {
		if item.ProductID == req.ProductID {
			existingQuantity = item.Quantity
			break
		}
	}

	// 檢查總數量是否超過庫存
	totalQuantity := existingQuantity + req.Quantity
	if totalQuantity > productInfo.Stock {
		return fmt.Errorf("total quantity exceeds stock: cart has %d, adding %d, stock is %d",
			existingQuantity, req.Quantity, productInfo.Stock)
	}

	// 獲取圖片 base64
	var imageBase64 string
	if len(productInfo.Images) > 0 {

		if productInfo.Images[0].Data != "" {
			imageBase64 = productInfo.Images[0].Data
			fmt.Printf("Using existing base64 image data with length: %d", len(imageBase64))
		} else {
			// 如果沒有 data 字段，則從 URL 獲取
			imageURL := productInfo.Images[0].URL
			if !strings.HasPrefix(imageURL, "http") {
				imageURL = fmt.Sprintf("%s%s", s.config.ProductServiceBaseURL, imageURL)
				log.Printf("Converted to absolute URL: %s", imageURL)
			}
			imageBase64, err = s.productClient.GetProductImageAsBase64(ctx, imageURL)
			if err != nil {
				log.Printf("Warning: Failed to get image as base64: %v", err)
				// 繼續執行，不因圖片獲取失敗而中斷流程
			} else {
				log.Printf("Successfully got base64 image with length: %d", len(imageBase64))
			}
		}
	} else {
		fmt.Printf("Product %s has no images", req.ProductID)
	}

	// 創建購物車項目
	item := model.CartItem{
		ProductID:  req.ProductID,
		Name:       productInfo.Name,
		Image:      imageBase64,
		Price:      productInfo.Price,
		Quantity:   req.Quantity,
		Selected:   true,
		StockCount: productInfo.Stock,
		UpdatedAt:  time.Now(),
	}

	log.Printf("Cart item created: %+v", item)

	// 添加到購物車
	if err := s.cartRepo.AddItem(ctx, userID, item); err != nil {
		log.Printf("Error adding item to cart: %v", err)
		return fmt.Errorf("failed to add item to cart: %w", err)
	}

	log.Printf("Successfully added product %s to cart for user %s", req.ProductID, userID)
	return nil
}

func (s *cartService) RemoveItem(ctx context.Context, userID string, productID string) error {
	return s.cartRepo.RemoveItem(ctx, userID, productID)
}

func (s *cartService) UpdateQuantity(ctx context.Context, userID string, req *model.UpdateQuantityRequest) error {
	// 調用 product service 檢查庫存
	productInfo, err := s.productClient.GetProduct(ctx, req.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get product info: %w", err)
	}
	if productInfo == nil {
		return ErrProductNotFound
	}

	if req.Quantity < 0 || req.Quantity > productInfo.Stock {
		return ErrInvalidStock
	}

	return s.cartRepo.UpdateQuantity(ctx, userID, req.ProductID, req.Quantity)
}

func (s *cartService) ClearCart(ctx context.Context, userID string) error {
	return s.cartRepo.ClearCart(ctx, userID)
}

func (s *cartService) SelectItems(ctx context.Context, userID string, req *model.SelectItemsRequest) error {
	// 檢查所有商品是否存在
	for _, productID := range req.ProductIDs {
		productInfo, err := s.productClient.GetProduct(ctx, productID)
		if err != nil {
			return fmt.Errorf("failed to get product info: %w", err)
		}
		if productInfo == nil {
			return fmt.Errorf("product not found: %s", productID)
		}
	}

	return s.cartRepo.SelectItems(ctx, userID, req.ProductIDs)
}

func (s *cartService) CreateOrder(ctx context.Context, userID string) error {
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("get cart failed: %w", err)
	}

	var selectedItems []client.CartItemInfo
	for _, item := range cart.Items {
		if item.Selected {
			selectedItems = append(selectedItems, client.CartItemInfo{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			})
		}
	}

	if len(selectedItems) == 0 {
		return fmt.Errorf("no items selected")
	}

	_, err = s.orderClient.CreateOrder(ctx, &client.CreateOrderRequest{
		UserID: userID,
		Items:  selectedItems,
	})
	if err != nil {
		return fmt.Errorf("create order failed: %w", err)
	}

	for _, item := range selectedItems {
		if err := s.cartRepo.RemoveItem(ctx, userID, item.ProductID); err != nil {
			return fmt.Errorf("remove item failed: %w", err)
		}
	}

	return nil
}

// 添加 GetCartItems 方法實現
func (s *cartService) GetCartItems(ctx context.Context, userID string) ([]model.CartItem, error) {
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	return cart.Items, nil
}
