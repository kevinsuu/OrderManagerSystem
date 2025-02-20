package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/order-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/order-service/internal/handler"
	"github.com/kevinsuu/OrderManagerSystem/order-service/internal/middleware"
	"github.com/kevinsuu/OrderManagerSystem/order-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/order-service/internal/service"
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
	orderRepo := repository.NewOrderRepository(db)
	
	// 初始化服務層
	orderService := service.NewOrderService(orderRepo, redisClient)

	// 初始化 HTTP 處理器
	handler := handler.NewHandler(orderService)

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
		// 添加 auth middleware
		api.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		
		orders := api.Group("/orders")
		{
			orders.POST("/", handler.CreateOrder)
			orders.GET("/", handler.ListOrders)
			orders.GET("/:id", handler.GetOrder)
			orders.PUT("/:id", handler.UpdateOrder)
			orders.DELETE("/:id", handler.DeleteOrder)
			orders.POST("/:id/cancel", handler.CancelOrder)
			orders.GET("/status/:status", handler.GetOrdersByStatus)
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
