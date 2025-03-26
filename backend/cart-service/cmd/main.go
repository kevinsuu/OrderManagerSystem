package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/client"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/handler"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/infrastructure/firebase"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/middleware"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/service"
)

// 新增這個函數
func checkServiceHealth(serviceURL, serviceName string) error {
	resp, err := http.Get(serviceURL + "/health")
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", serviceName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s health check failed with status: %d", serviceName, resp.StatusCode)
	}

	log.Printf("%s is healthy", serviceName)
	return nil
}

func main() {
	// 加載配置
	cfg := config.LoadConfig()
	if err := checkServiceHealth(cfg.ProductService.BaseURL, "Product Service"); err != nil {
		fmt.Printf("Warning: %v", err)
	}

	if err := checkServiceHealth(cfg.OrderService.BaseURL, "Order Service"); err != nil {
		fmt.Printf("Warning: %v", err)
	}

	// 初始化 Firebase
	ctx := context.Background()
	fb, err := firebase.InitFirebase(ctx, cfg.Firebase.CredentialsFile, cfg.Firebase.ProjectID)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// 初始化存儲層
	cartRepo := repository.NewCartRepository(fb.Database)

	// 初始化客戶端
	productClient := client.NewProductClient(cfg.ProductService.BaseURL)
	orderClient := client.NewOrderClient(cfg.OrderService.BaseURL)

	// 初始化服務層
	cartService := service.NewCartService(cartRepo, productClient, orderClient)

	// 初始化 HTTP 處理器
	handler := handler.NewHandler(cartService)

	// 設置 Gin 路由
	router := gin.Default()

	// 中間件
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// CORS 中間件配置
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://your-production-domain.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
			cart.POST("/checkout", handler.CreateOrder)
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
