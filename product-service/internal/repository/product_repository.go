package repository

import (
	"context"
	"errors"

	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/model"
	"gorm.io/gorm"
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
	db *gorm.DB
}

// NewProductRepository 創建產品存儲實例
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create 創建產品
func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID 根據ID獲取產品
func (r *productRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).
		Preload("Images").
		Preload("Attributes").
		First(&product, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// Update 更新產品
func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新產品基本信息
		if err := tx.Save(product).Error; err != nil {
			return err
		}

		// 更新圖片
		if len(product.Images) > 0 {
			if err := tx.Where("product_id = ?", product.ID).Delete(&model.Image{}).Error; err != nil {
				return err
			}
			if err := tx.Create(&product.Images).Error; err != nil {
				return err
			}
		}

		// 更新屬性
		if len(product.Attributes) > 0 {
			if err := tx.Where("product_id = ?", product.ID).Delete(&model.Attribute{}).Error; err != nil {
				return err
			}
			if err := tx.Create(&product.Attributes).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete 刪除產品
func (r *productRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 刪除相關的圖片和屬性
		if err := tx.Where("product_id = ?", id).Delete(&model.Image{}).Error; err != nil {
			return err
		}
		if err := tx.Where("product_id = ?", id).Delete(&model.Attribute{}).Error; err != nil {
			return err
		}
		// 軟刪除產品
		return tx.Delete(&model.Product{}, "id = ?", id).Error
	})
}

// List 獲取產品列表
func (r *productRepository) List(ctx context.Context, page, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Images").
		Preload("Attributes").
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetByCategoryID 根據分類獲取產品
func (r *productRepository) GetByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("category_id = ?", categoryID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Images").
		Preload("Attributes").
		Where("category_id = ?", categoryID).
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// UpdateStock 更新庫存
func (r *productRepository) UpdateStock(ctx context.Context, id string, quantity int) error {
	result := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ? AND stock >= ?", id, -quantity).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}

	return nil
}

// SearchProducts 搜索產品
func (r *productRepository) SearchProducts(ctx context.Context, query string, page, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	offset := (page - 1) * limit

	searchQuery := "%" + query + "%"

	if err := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("name LIKE ? OR description LIKE ?", searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Images").
		Preload("Attributes").
		Where("name LIKE ? OR description LIKE ?", searchQuery, searchQuery).
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
