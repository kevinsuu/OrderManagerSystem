package service

import (
	"context"
	"fmt"
	"log"

	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/client"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/repository"
)

// WishlistService 收藏清單服務接口
type WishlistService interface {
	AddToWishlist(ctx context.Context, userId, productId string) error
	RemoveFromWishlist(ctx context.Context, userId, productId string) error
	GetWishlist(ctx context.Context, userId string, page, limit int) (*model.WishlistResponse, error)
	GetWishlistWithProductDetails(ctx context.Context, userId string, page, limit int) (*model.WishlistResponse, error)
	IsProductInWishlist(ctx context.Context, userId, productId string) (bool, error)
}

// wishlistService 實現 WishlistService 接口
type wishlistService struct {
	wishlistRepo  repository.WishlistRepository
	productClient client.ProductClient
}

// NewWishlistService 創建一個新的收藏清單服務
func NewWishlistService(wishlistRepo repository.WishlistRepository, productClient client.ProductClient) WishlistService {
	return &wishlistService{
		wishlistRepo:  wishlistRepo,
		productClient: productClient,
	}
}

// AddToWishlist 添加商品到收藏清單
func (s *wishlistService) AddToWishlist(ctx context.Context, userId, productId string) error {
	// 檢查商品是否存在
	_, err := s.productClient.GetProductById(ctx, productId)
	if err != nil {
		log.Printf("Error checking product %s: %v", productId, err)
		return fmt.Errorf("product not found: %w", err)
	}

	// 先檢查商品是否已在收藏清單中
	exists, err := s.wishlistRepo.IsProductInWishlist(ctx, userId, productId)
	if err != nil {
		return fmt.Errorf("failed to check wishlist: %w", err)
	}
	if exists {
		return fmt.Errorf("product already in wishlist")
	}

	// 添加到收藏清單
	return s.wishlistRepo.AddToWishlist(ctx, userId, productId)
}

// RemoveFromWishlist 從收藏清單移除商品
func (s *wishlistService) RemoveFromWishlist(ctx context.Context, userId, productId string) error {
	return s.wishlistRepo.RemoveFromWishlist(ctx, userId, productId)
}

// GetWishlist 獲取使用者的收藏清單
func (s *wishlistService) GetWishlist(ctx context.Context, userId string, page, limit int) (*model.WishlistResponse, error) {
	return s.wishlistRepo.GetWishlist(ctx, userId, page, limit)
}

// GetWishlistWithProductDetails 獲取使用者的收藏清單，並包含商品詳細資訊
func (s *wishlistService) GetWishlistWithProductDetails(ctx context.Context, userId string, page, limit int) (*model.WishlistResponse, error) {
	// 獲取收藏清單
	wishlistResp, err := s.wishlistRepo.GetWishlist(ctx, userId, page, limit)
	if err != nil {
		return nil, err
	}

	// 如果清單為空，直接返回
	if len(wishlistResp.Wishlist) == 0 {
		return wishlistResp, nil
	}

	// 豐富商品詳細資訊
	for i, item := range wishlistResp.Wishlist {
		product, err := s.productClient.GetProductById(ctx, item.ProductId)
		if err != nil {
			log.Printf("Error getting product details for %s: %v", item.ProductId, err)
			continue
		}

		// 在這裡可以設置更多商品詳細資訊
		wishlistResp.Wishlist[i].Product = product
	}

	return wishlistResp, nil
}

// IsProductInWishlist 檢查商品是否在收藏清單中
func (s *wishlistService) IsProductInWishlist(ctx context.Context, userId, productId string) (bool, error) {
	return s.wishlistRepo.IsProductInWishlist(ctx, userId, productId)
}
