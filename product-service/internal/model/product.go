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
	ID          string        `json:"id" gorm:"primaryKey"`
	Name        string        `json:"name" gorm:"index"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Stock       int           `json:"stock"`
	Status      ProductStatus `json:"status"`
	CategoryID  string        `json:"categoryId" gorm:"index"`
	Images      []Image       `json:"images" gorm:"foreignKey:ProductID"`
	Attributes  []Attribute   `json:"attributes" gorm:"foreignKey:ProductID"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	DeletedAt   *time.Time    `json:"deletedAt,omitempty" gorm:"index"`
}

// Image 產品圖片
type Image struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ProductID string    `json:"productId" gorm:"index"`
	URL       string    `json:"url"`
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
	CategoryID  *string        `json:"categoryId,omitempty"`
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
