package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"firebase.google.com/go/db"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/model"
)

// ProductRepository 產品存儲接口
type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id string) (*model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]model.Product, int64, error)
	GetByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]model.Product, int64, error)
	UpdateStock(ctx context.Context, id string, quantity int) error
	SearchProducts(ctx context.Context, query string, page, limit int) ([]model.Product, int64, error)
}

type productRepository struct {
	db *db.Client
}

// NewProductRepository 創建產品存儲實例
func NewProductRepository(db *db.Client) *productRepository {
	return &productRepository{db: db}
}

// Create 創建產品
func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	ref := r.db.NewRef("products")
	newRef, err := ref.Push(ctx, nil)
	if err != nil {
		return fmt.Errorf("error creating product: %v", err)
	}

	product.ID = newRef.Key
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	if err := newRef.Set(ctx, product); err != nil {
		return fmt.Errorf("error saving product: %v", err)
	}

	return nil
}

// GetByID 根據ID獲取產品
func (r *productRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	ref := r.db.NewRef("products").Child(id)
	var product model.Product
	if err := ref.Get(ctx, &product); err != nil {
		if err.Error() == "http error status: 404; reason: Permission denied" {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting product: %v", err)
	}

	// 檢查是否為空數據
	if product.ID == "" {
		return nil, nil
	}

	// 確保 ID 被正確設置
	product.ID = id
	return &product, nil
}

// List 獲取產品列表
func (r *productRepository) List(ctx context.Context, page, limit int) ([]model.Product, int64, error) {
	ref := r.db.NewRef("products")
	var products map[string]model.Product
	if err := ref.Get(ctx, &products); err != nil {
		return nil, 0, fmt.Errorf("error getting products: %v", err)
	}

	result := make([]model.Product, 0, len(products))
	for _, product := range products {
		result = append(result, product)
	}

	total := int64(len(result))
	start := (page - 1) * limit
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	if start >= len(result) {
		return []model.Product{}, total, nil
	}

	return result[start:end], total, nil
}

// Update 更新產品
func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	ref := r.db.NewRef("products").Child(product.ID)
	product.UpdatedAt = time.Now()
	if err := ref.Set(ctx, product); err != nil {
		return fmt.Errorf("error updating product: %v", err)
	}
	return nil
}

// Delete 刪除產品
func (r *productRepository) Delete(ctx context.Context, id string) error {
	ref := r.db.NewRef("products").Child(id)
	if err := ref.Delete(ctx); err != nil {
		return fmt.Errorf("error deleting product: %v", err)
	}
	return nil
}

// GetByCategoryID 根據分類獲取產品
func (r *productRepository) GetByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]model.Product, int64, error) {
	ref := r.db.NewRef("products")
	var products map[string]model.Product
	if err := ref.Get(ctx, &products); err != nil {
		return nil, 0, fmt.Errorf("error getting products: %v", err)
	}

	var filteredProducts []model.Product
	for _, product := range products {
		if product.Category == categoryID {
			filteredProducts = append(filteredProducts, product)
		}
	}

	total := int64(len(filteredProducts))
	start := (page - 1) * limit
	end := start + limit
	if end > len(filteredProducts) {
		end = len(filteredProducts)
	}

	return filteredProducts[start:end], total, nil
}

// UpdateStock 更新庫存
func (r *productRepository) UpdateStock(ctx context.Context, id string, quantity int) error {
	ref := r.db.NewRef("products").Child(id)
	var product model.Product
	if err := ref.Get(ctx, &product); err != nil {
		return fmt.Errorf("error getting product: %v", err)
	}

	if product.Stock+quantity < 0 {
		return fmt.Errorf("insufficient stock")
	}

	product.Stock += quantity
	product.UpdatedAt = time.Now()

	if err := ref.Set(ctx, &product); err != nil {
		return fmt.Errorf("error updating stock: %v", err)
	}

	return nil
}

// SearchProducts 搜索產品
func (r *productRepository) SearchProducts(ctx context.Context, query string, page, limit int) ([]model.Product, int64, error) {
	ref := r.db.NewRef("products")
	var products map[string]model.Product
	if err := ref.Get(ctx, &products); err != nil {
		return nil, 0, fmt.Errorf("error getting products: %v", err)
	}

	var filteredProducts []model.Product
	for _, product := range products {
		// 檢查產品名稱或描述是否包含搜索關鍵字
		if containsIgnoreCase(product.Name, query) || containsIgnoreCase(product.Description, query) {
			filteredProducts = append(filteredProducts, product)
			continue
		}

		// 檢查產品類別是否匹配
		if product.Category != "" {
			categoryRef := r.db.NewRef("categories").Child(product.Category)
			var category model.Category
			if err := categoryRef.Get(ctx, &category); err == nil && category.ID != "" {
				if containsIgnoreCase(category.Name, query) {
					filteredProducts = append(filteredProducts, product)
				}
			}
		}
	}

	total := int64(len(filteredProducts))
	start := (page - 1) * limit
	end := start + limit
	if end > len(filteredProducts) {
		end = len(filteredProducts)
	}
	if start >= len(filteredProducts) {
		return []model.Product{}, total, nil
	}

	return filteredProducts[start:end], total, nil
}

func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}
