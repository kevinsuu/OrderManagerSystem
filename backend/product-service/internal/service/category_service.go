package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/model"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/repository"
	"time"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req *model.CreateCategoryRequest) (*model.Category, error)
	GetCategory(ctx context.Context, id string) (*model.Category, error)
	UpdateCategory(ctx context.Context, id string, req *model.UpdateCategoryRequest) (*model.Category, error)
	DeleteCategory(ctx context.Context, id string) error
	ListCategories(ctx context.Context) ([]model.Category, error)
	GetSubcategories(ctx context.Context, parentID string) ([]model.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(ctx context.Context, req *model.CreateCategoryRequest) (*model.Category, error) {
	category := &model.Category{
		ID:        uuid.New().String(),
		Name:      req.Name,
		ParentID:  req.ParentID,
		Level:     req.Level,
		Sort:      req.Sort,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetCategory(ctx context.Context, id string) (*model.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id string, req *model.UpdateCategoryRequest) (*model.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.Level != nil {
		category.Level = *req.Level
	}
	if req.Sort != nil {
		category.Sort = *req.Sort
	}
	category.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return ErrCategoryNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *categoryService) ListCategories(ctx context.Context) ([]model.Category, error) {
	return s.repo.List(ctx)
}

func (s *categoryService) GetSubcategories(ctx context.Context, parentID string) ([]model.Category, error) {
	return s.repo.GetSubcategories(ctx, parentID)
}

// ... 實現其他服務方法 ... 