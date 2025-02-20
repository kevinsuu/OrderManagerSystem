package handler

import (
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/service"
)

// Handler 處理所有HTTP請求的結構體
type Handler struct {
	productService  service.ProductService
	categoryService service.CategoryService
}

// NewHandler 創建新的Handler實例
func NewHandler(productService service.ProductService, categoryService service.CategoryService) *Handler {
	return &Handler{
		productService:  productService,
		categoryService: categoryService,
	}
}
