package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/client"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/handler"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/middleware"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/service"
)

func main() {
	// 加載配置
	cfg := config.LoadConfig()

	// 初始化資料庫連接
	db := repository.NewPostgresDB(cfg.Database)

	// 初始化 Redis
	redisClient := repository.NewRedisRepository(cfg.Redis)
	defer redisClient.Close()

	// 初始化 product client
	productClient := client.NewProductClient(cfg.ProductService.BaseURL)

	// 初始化存儲層
	cartRepo := repository.NewCartRepository(redisClient, db)

	// 初始化服務層
	cartService := service.NewCartService(cartRepo, productClient)

	// 初始化 HTTP 處理器
	handler := handler.NewHandler(cartService)

	// 設置 Gin 路由
	router := gin.Default()

	// 中間件
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// 健康檢查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由
	api := router.Group("/api/v1")
	{
		// 添加 auth middleware
		api.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

		cart := api.Group("/cart")
		{
			cart.GET("/", handler.GetCart)
			cart.POST("/items", handler.AddToCart)
			cart.DELETE("/items/:productId", handler.RemoveFromCart)
			cart.PUT("/items", handler.UpdateQuantity)
			cart.POST("/items/select", handler.SelectItems)
			cart.DELETE("/", handler.ClearCart)
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
