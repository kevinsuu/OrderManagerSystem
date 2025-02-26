package service

import (
	"context"
	"errors"

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
	cartRepo repository.CartRepository
	// TODO: 添加 product service client 用於檢查商品信息和庫存
}

func NewCartService(cartRepo repository.CartRepository) CartService {
	return &cartService{
		cartRepo: cartRepo,
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

// ... 實現其他方法 ...
