package repository

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/db"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/model"
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
	db *db.Client
}

func NewCategoryRepository(db *db.Client) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	ref := r.db.NewRef("categories")
	newRef, err := ref.Push(ctx, nil)
	if err != nil {
		return fmt.Errorf("error creating category: %v", err)
	}

	category.ID = newRef.Key
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	if err := newRef.Set(ctx, category); err != nil {
		return fmt.Errorf("error saving category: %v", err)
	}

	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id string) (*model.Category, error) {
	ref := r.db.NewRef("categories").Child(id)
	var category model.Category
	if err := ref.Get(ctx, &category); err != nil {
		if err.Error() == "http error status: 404; reason: Permission denied" {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting category: %v", err)
	}

	if category.ID == "" {
		return nil, nil
	}

	category.ID = id
	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	ref := r.db.NewRef("categories").Child(category.ID)
	category.UpdatedAt = time.Now()
	if err := ref.Set(ctx, category); err != nil {
		return fmt.Errorf("error updating category: %v", err)
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	ref := r.db.NewRef("categories").Child(id)
	if err := ref.Delete(ctx); err != nil {
		return fmt.Errorf("error deleting category: %v", err)
	}
	return nil
}

func (r *categoryRepository) List(ctx context.Context) ([]model.Category, error) {
	ref := r.db.NewRef("categories")
	var categories map[string]model.Category
	if err := ref.Get(ctx, &categories); err != nil {
		return nil, fmt.Errorf("error getting categories: %v", err)
	}

	result := make([]model.Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, category)
	}
	return result, nil
}

func (r *categoryRepository) GetSubcategories(ctx context.Context, parentID string) ([]model.Category, error) {
	ref := r.db.NewRef("categories")
	var categories map[string]model.Category
	if err := ref.Get(ctx, &categories); err != nil {
		return nil, fmt.Errorf("error getting categories: %v", err)
	}

	var result []model.Category
	for _, category := range categories {
		if category.ParentID != nil && *category.ParentID == parentID {
			result = append(result, category)
		}
	}
	return result, nil
}
