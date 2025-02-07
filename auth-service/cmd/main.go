package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/handler"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/service"
)

func main() {
	// 加載配置
	cfg := config.LoadConfig()

	// 初始化資料庫連接
	db := repository.NewPostgresDB(cfg.Database)

	// 初始化存儲層
	userRepo := repository.NewUserRepository(db)

	// 初始化服務層
	authService := service.NewAuthService(userRepo, "your-secret-key", 24*time.Hour)

	// 初始化 HTTP 處理器
	handler := handler.NewHandler(authService)

	// 設置 Gin 路由
	router := gin.Default()
	
	// 中間件
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// 路由組
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
			auth.POST("/refresh", handler.RefreshToken)
		}

		// 健康檢查
		api.GET("/health", handler.HealthCheck)
		
		// 需要認證的路由
		secured := api.Group("/")
		{
			secured.GET("/validate", handler.ValidateToken)
			secured.GET("/user/:id", handler.GetUser)
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
