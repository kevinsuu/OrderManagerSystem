package repository

import (
	"context"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	GetByID(ctx context.Context, id string) (*model.Category, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.Category, error)
	GetSubcategories(ctx context.Context, parentID string) ([]model.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) GetByID(ctx context.Context, id string) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, "id = ?", id).Error
}

func (r *categoryRepository) List(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) GetSubcategories(ctx context.Context, parentID string) ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
} 