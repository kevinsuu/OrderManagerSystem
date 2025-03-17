package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/model"
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

// HealthCheck 健康檢查
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreateProduct 創建產品
func (h *Handler) CreateProduct(c *gin.Context) {
	var req model.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct 獲取產品
func (h *Handler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "獲取產品失敗",
		})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "product not found",
			"message": "找不到該產品",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "獲取產品成功",
		"data":    product,
	})
}

// ListProducts 獲取產品列表
func (h *Handler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.productService.List(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "獲取產品列表失敗",
		})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "暫無產品數據",
			"data": gin.H{
				"products": []interface{}{},
				"total":    0,
				"page":     page,
				"limit":    limit,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "獲取產品列表成功",
		"data": gin.H{
			"products": products,
			"total":    total,
			"page":     page,
			"limit":    limit,
		},
	})
}

// UpdateProduct 更新產品
func (h *Handler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct 刪除產品
func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "product not found",
				"message": "找不到該產品",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "刪除產品失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "產品刪除成功",
	})
}

// GetProductsByCategory 獲取分類下的產品
func (h *Handler) GetProductsByCategory(c *gin.Context) {
	categoryID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.productService.GetByCategoryID(c.Request.Context(), categoryID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// UpdateStock 更新庫存
func (h *Handler) UpdateStock(c *gin.Context) {
	id := c.Param("id")
	var req model.StockUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "請求格式錯誤",
		})
		return
	}

	if err := h.productService.UpdateStock(c.Request.Context(), id, req.Quantity); err != nil {
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "product not found",
				"message": "找不到該產品",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "更新庫存失敗",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "庫存更新成功",
	})
}

// SearchProducts 搜索產品
func (h *Handler) SearchProducts(c *gin.Context) {
	query := c.Query("query")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.productService.SearchProducts(c.Request.Context(), query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "搜索產品失敗",
		})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "未找到相關產品",
			"data": gin.H{
				"products": []interface{}{},
				"total":    0,
				"page":     page,
				"limit":    limit,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "搜索產品成功",
		"data": gin.H{
			"products": products,
			"total":    total,
			"page":     page,
			"limit":    limit,
		},
	})
}
