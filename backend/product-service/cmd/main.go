package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/handler"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/infrastructure/firebase"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/middleware"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/product-service/internal/service"
)

func main() {
	// 加載配置
	cfg := config.LoadConfig()

	// 初始化 Firebase
	ctx := context.Background()
	fb, err := firebase.InitFirebase(ctx, cfg.Firebase.CredentialsFile, cfg.Firebase.ProjectID)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// 初始化存儲層
	productRepo := repository.NewProductRepository(fb.Database)
	categoryRepo := repository.NewCategoryRepository(fb.Database)

	// 初始化服務層
	productService := service.NewProductService(productRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// 初始化 HTTP 處理器
	handler := handler.NewHandler(productService, categoryService)

	// 設置 Gin 路由
	router := gin.Default()

	// 中間件
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// CORS 中間件配置
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	// 路由組
	api := router.Group("/api/v1")

	// 公開路由（不需要驗證）
	{
		// 產品查詢相關路由
		products := api.Group("/products")
		{
			products.GET("/", handler.ListProducts)
			products.GET("/search", handler.SearchProducts)
			products.GET("/:id", handler.GetProduct)
			products.GET("/category/:id", handler.GetProductsByCategory)
		}

		// 分類查詢相關路由
		categories := api.Group("/categories")
		{
			categories.GET("/", handler.ListCategories)
			categories.GET("/:id", handler.GetCategory)
			categories.GET("/:id/products", handler.GetProductsByCategory)
		}
	}

	// 需要驗證的路由
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		// 產品管理路由
		products := protected.Group("/products")
		{
			products.POST("/", handler.CreateProduct)
			products.PUT("/:id", handler.UpdateProduct)
			products.PUT("/:id/stock", handler.UpdateStock)
			products.DELETE("/:id", handler.DeleteProduct)
		}

		// 分類管理路由
		categories := protected.Group("/categories")
		{
			categories.POST("/", handler.CreateCategory)
			categories.PUT("/:id", handler.UpdateCategory)
			categories.DELETE("/:id", handler.DeleteCategory)
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
