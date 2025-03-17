package service

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	Create(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error)
	GetByID(ctx context.Context, id string) (*model.Product, error)
	List(ctx context.Context, page, limit int) ([]model.Product, int64, error)
	Update(ctx context.Context, id string, req *model.UpdateProductRequest) (*model.Product, error)
	Delete(ctx context.Context, id string) error
	GetByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]model.Product, int64, error)
	UpdateStock(ctx context.Context, id string, quantity int) error
	SearchProducts(ctx context.Context, query string, page, limit int) ([]model.Product, int64, error)
}

type productService struct {
	repo repository.ProductRepository
}

// NewProductService 創建產品服務實例
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

// Create 創建產品
func (s *productService) Create(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Status:      model.ProductStatusActive,
		Category:    req.CategoryID,
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

	return product, nil
}

// GetByID 獲取產品
func (s *productService) GetByID(ctx context.Context, id string) (*model.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// List 獲取產品列表
func (s *productService) List(ctx context.Context, page, limit int) ([]model.Product, int64, error) {
	return s.repo.List(ctx, page, limit)
}

// Update 更新產品
func (s *productService) Update(ctx context.Context, id string, req *model.UpdateProductRequest) (*model.Product, error) {
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
	if req.Category != nil {
		product.Category = *req.Category
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

	return product, nil
}

// Delete 刪除產品
func (s *productService) Delete(ctx context.Context, id string) error {
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

	return nil
}

// GetByCategoryID 獲取分類產品
func (s *productService) GetByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]model.Product, int64, error) {
	return s.repo.GetByCategoryID(ctx, categoryID, page, limit)
}

// UpdateStock 更新庫存
func (s *productService) UpdateStock(ctx context.Context, id string, quantity int) error {
	return s.repo.UpdateStock(ctx, id, quantity)
}

// SearchProducts 搜索產品
func (s *productService) SearchProducts(ctx context.Context, query string, page, limit int) ([]model.Product, int64, error) {
	return s.repo.SearchProducts(ctx, query, page, limit)
}
