package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/handler"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/middleware"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/service"
)

func main() {
	// 加載配置
	cfg := config.LoadConfig()

	// 初始化資料庫連接
	db := repository.NewPostgresDB(cfg.Database)

	// 初始化 Redis (用於快取)
	redisClient := repository.NewRedisRepository(cfg.Redis)
	defer redisClient.Close()

	// 初始化存儲層
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// 初始化服務層
	productService := service.NewProductService(productRepo, redisClient)
	categoryService := service.NewCategoryService(categoryRepo)

	// 初始化 HTTP 處理器
	handler := handler.NewHandler(productService, categoryService)

	// 設置 Gin 路由
	router := gin.Default()

	// 基本中間件
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// 健康檢查 (不需要驗證)
	router.GET("/health", handler.HealthCheck)

	// API 路由 (需要驗證)
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		products := api.Group("/products")
		{
			products.POST("/", handler.CreateProduct)
			products.GET("/", handler.ListProducts)
			products.GET("/:id", handler.GetProduct)
			products.PUT("/:id", handler.UpdateProduct)
			products.DELETE("/:id", handler.DeleteProduct)
			products.PUT("/:id/stock", handler.UpdateStock)
			products.GET("/category/:categoryId", handler.GetProductsByCategory)
			products.GET("/search", handler.SearchProducts)
		}

		categories := api.Group("/categories")
		{
			categories.POST("/", handler.CreateCategory)
			categories.GET("/", handler.ListCategories)
			categories.GET("/:id", handler.GetCategory)
			categories.PUT("/:id", handler.UpdateCategory)
			categories.DELETE("/:id", handler.DeleteCategory)
			categories.GET("/:id/subcategories", handler.GetSubcategories)
		}
	}

	// 啟動服務器
	go func() {
		if err := router.Run(cfg.Server.Address); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 優雅關閉
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
