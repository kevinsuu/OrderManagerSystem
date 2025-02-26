package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/client"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/repository"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidStock    = errors.New("invalid stock quantity")
)

type CartService interface {
	GetCart(ctx context.Context, userID string) (*model.CartResponse, error)
	AddItem(ctx context.Context, userID string, req *model.AddToCartRequest) error
	RemoveItem(ctx context.Context, userID string, productID string) error
	UpdateQuantity(ctx context.Context, userID string, req *model.UpdateQuantityRequest) error
	ClearCart(ctx context.Context, userID string) error
	SelectItems(ctx context.Context, userID string, req *model.SelectItemsRequest) error
}

type cartService struct {
	cartRepo      repository.CartRepository
	productClient client.ProductClient
}

func NewCartService(cartRepo repository.CartRepository, productClient client.ProductClient) CartService {
	return &cartService{
		cartRepo:      cartRepo,
		productClient: productClient,
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

func (s *cartService) AddItem(ctx context.Context, userID string, req *model.AddToCartRequest) error {
	// 調用 product service 檢查商品是否存在及庫存
	productInfo, err := s.productClient.GetProduct(ctx, req.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get product info: %w", err)
	}
	if productInfo == nil {
		return ErrProductNotFound
	}

	// 檢查庫存
	if productInfo.StockCount < req.Quantity {
		return ErrInvalidStock
	}

	item := model.CartItem{
		ProductID:  req.ProductID,
		Name:       productInfo.Name,
		Image:      productInfo.Image,
		Price:      productInfo.Price,
		Quantity:   req.Quantity,
		Selected:   true,
		StockCount: productInfo.StockCount,
		UpdatedAt:  time.Now(),
	}

	return s.cartRepo.AddItem(ctx, userID, item)
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

	if req.Quantity < 0 || req.Quantity > productInfo.StockCount {
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
