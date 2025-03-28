package model

import (
	"time"
)

// ProductStatus 產品狀態
type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusSoldOut  ProductStatus = "sold_out"
)

// Product 產品模型
type Product struct {
	ID          string        `json:"id" db:"id"`
	Name        string        `json:"name" db:"name"`
	Description string        `json:"description" db:"description"`
	Price       float64       `json:"price" db:"price"`
	Stock       int           `json:"stock" db:"stock"`
	Status      ProductStatus `json:"status" db:"status"`
	Category    string        `json:"category" db:"category"`
	Images      []Image       `json:"images" gorm:"foreignKey:ProductID"`
	Attributes  []Attribute   `json:"attributes,omitempty" gorm:"foreignKey:ProductID"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time    `json:"deleted_at,omitempty" gorm:"index"`
}

// Image 產品圖片
type Image struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ProductID string    `json:"productId" gorm:"index"`
	Data      string    `json:"data"` // base64 圖片數據
	URL       string    `json:"url"`  // 圖片 URL
	Sort      int       `json:"sort"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Attribute 產品屬性
type Attribute struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ProductID string    `json:"productId" gorm:"index"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateProductRequest 創建產品請求
type CreateProductRequest struct {
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description" binding:"required"`
	Price       float64     `json:"price" binding:"required,gt=0"`
	Stock       int         `json:"stock" binding:"required,gte=0"`
	CategoryID  string      `json:"categoryId" binding:"required"`
	Images      []Image     `json:"images" binding:"required,min=1"`
	Attributes  []Attribute `json:"attributes"`
}

// UpdateProductRequest 更新產品請求
type UpdateProductRequest struct {
	Name        *string        `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	Price       *float64       `json:"price,omitempty" binding:"omitempty,gt=0"`
	Stock       *int           `json:"stock,omitempty" binding:"omitempty,gte=0"`
	Status      *ProductStatus `json:"status,omitempty"`
	Category    *string        `json:"category,omitempty"`
	Images      []Image        `json:"images,omitempty"`
	Attributes  []Attribute    `json:"attributes,omitempty"`
}

// ProductResponse 產品響應
type ProductResponse struct {
	Product
	Category *Category `json:"category,omitempty"`
}

// ProductListResponse 產品列表響應
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

// StockUpdateRequest 庫存更新請求
type StockUpdateRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

// ProductRequest 產品請求
type ProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	Category    string  `json:"category" binding:"required"`
}
