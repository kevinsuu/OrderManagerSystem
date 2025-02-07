package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/handler"
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

	// 初始化服務層
	productService := service.NewProductService(productRepo, redisClient)

	// 初始化 HTTP 處理器
	handler := handler.NewHandler(productService)

	// 設置 Gin 路由
	router := gin.Default()

	// 中間件
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// 健康檢查
	router.GET("/health", handler.HealthCheck)

	// API 路由
	api := router.Group("/api/v1")
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
