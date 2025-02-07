package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/repository"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidStock    = errors.New("invalid stock quantity")
)

// ProductService 產品服務接口
type ProductService interface {
	CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error)
	GetProduct(ctx context.Context, id string) (*model.ProductResponse, error)
	UpdateProduct(ctx context.Context, id string, req *model.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, page, limit int) (*model.ProductListResponse, error)
	GetProductsByCategory(ctx context.Context, categoryID string, page, limit int) (*model.ProductListResponse, error)
	UpdateStock(ctx context.Context, id string, req *model.StockUpdateRequest) error
	SearchProducts(ctx context.Context, query string, page, limit int) (*model.ProductListResponse, error)
}

type productService struct {
	repo  repository.ProductRepository
	redis *redis.Client
}

// NewProductService 創建產品服務實例
func NewProductService(repo repository.ProductRepository, redis *redis.Client) ProductService {
	return &productService{
		repo:  repo,
		redis: redis,
	}
}

// CreateProduct 創建產品
func (s *productService) CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Status:      model.ProductStatusActive,
		CategoryID:  req.CategoryID,
		Images:      req.Images,
		Attributes:  req.Attributes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 為每個圖片生成ID
	for i := range product.Images {
		product.Images[i].ID = uuid.New().String()
		product.Images[i].ProductID = product.ID
		product.Images[i].CreatedAt = time.Now()
		product.Images[i].UpdatedAt = time.Now()
	}

	// 為每個屬性生成ID
	for i := range product.Attributes {
		product.Attributes[i].ID = uuid.New().String()
		product.Attributes[i].ProductID = product.ID
		product.Attributes[i].CreatedAt = time.Now()
		product.Attributes[i].UpdatedAt = time.Now()
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// 清除相關快取
	s.clearProductCache(ctx, product.CategoryID)

	return product, nil
}

// GetProduct 獲取產品
func (s *productService) GetProduct(ctx context.Context, id string) (*model.ProductResponse, error) {
	// 嘗試從快取獲取
	cacheKey := fmt.Sprintf("product:%s", id)
	if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
		var response model.ProductResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	response := &model.ProductResponse{
		Product: *product,
	}

	// 設置快取
	if cached, err := json.Marshal(response); err == nil {
		s.redis.Set(ctx, cacheKey, cached, 1*time.Hour)
	}

	return response, nil
}

// UpdateProduct 更新產品
func (s *productService) UpdateProduct(ctx context.Context, id string, req *model.UpdateProductRequest) (*model.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// 更新產品信息
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if len(req.Images) > 0 {
		product.Images = req.Images
		for i := range product.Images {
			if product.Images[i].ID == "" {
				product.Images[i].ID = uuid.New().String()
			}
			product.Images[i].ProductID = product.ID
			product.Images[i].UpdatedAt = time.Now()
		}
	}
	if len(req.Attributes) > 0 {
		product.Attributes = req.Attributes
		for i := range product.Attributes {
			if product.Attributes[i].ID == "" {
				product.Attributes[i].ID = uuid.New().String()
			}
			product.Attributes[i].ProductID = product.ID
			product.Attributes[i].UpdatedAt = time.Now()
		}
	}
	product.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// 清除快取
	s.clearProductCache(ctx, product.CategoryID)
	s.redis.Del(ctx, fmt.Sprintf("product:%s", id))

	return product, nil
}

// DeleteProduct 刪除產品
func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}
	if product == nil {
		return ErrProductNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// 清除快取
	s.clearProductCache(ctx, product.CategoryID)
	s.redis.Del(ctx, fmt.Sprintf("product:%s", id))

	return nil
}

// ListProducts 獲取產品列表
func (s *productService) ListProducts(ctx context.Context, page, limit int) (*model.ProductListResponse, error) {
	products, total, err := s.repo.List(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	response := &model.ProductListResponse{
		Products: make([]model.ProductResponse, len(products)),
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	for i, product := range products {
		response.Products[i] = model.ProductResponse{Product: product}
	}

	return response, nil
}

// GetProductsByCategory 獲取分類產品
func (s *productService) GetProductsByCategory(ctx context.Context, categoryID string, page, limit int) (*model.ProductListResponse, error) {
	products, total, err := s.repo.GetByCategoryID(ctx, categoryID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}

	response := &model.ProductListResponse{
		Products: make([]model.ProductResponse, len(products)),
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	for i, product := range products {
		response.Products[i] = model.ProductResponse{Product: product}
	}

	return response, nil
}

// UpdateStock 更新庫存
func (s *productService) UpdateStock(ctx context.Context, id string, req *model.StockUpdateRequest) error {
	if err := s.repo.UpdateStock(ctx, id, req.Quantity); err != nil {
		if err.Error() == "insufficient stock" {
			return ErrInvalidStock
		}
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// 清除快取
	s.redis.Del(ctx, fmt.Sprintf("product:%s", id))

	return nil
}

// SearchProducts 搜索產品
func (s *productService) SearchProducts(ctx context.Context, query string, page, limit int) (*model.ProductListResponse, error) {
	products, total, err := s.repo.SearchProducts(ctx, query, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	response := &model.ProductListResponse{
		Products: make([]model.ProductResponse, len(products)),
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	for i, product := range products {
		response.Products[i] = model.ProductResponse{Product: product}
	}

	return response, nil
}

// clearProductCache 清除產品相關快取
func (s *productService) clearProductCache(ctx context.Context, categoryID string) {
	s.redis.Del(ctx, fmt.Sprintf("category:%s:products", categoryID))
}
