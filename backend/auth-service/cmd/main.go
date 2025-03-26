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
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/handler"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/infrastructure/firebase"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/middleware"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/service"
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
	userRepo := repository.NewUserRepository(fb.Database)

	// 初始化服務層
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, time.Duration(cfg.JWT.ExpiryMinutes)*time.Minute)

	// 初始化 HTTP 處理器
	handler := handler.NewHandler(authService)

	// 設置 Gin 路由
	router := gin.New()

	// CORS 中間件配置
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 添加基本中間件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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
		secured.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			secured.GET("/validate", handler.ValidateToken)
			secured.GET("/user/:id", handler.GetUser)

			// 用戶相關路由組
			user := secured.Group("/user")
			{
				// 地址管理
				addresses := user.Group("/addresses")
				{
					addresses.POST("/", handler.CreateAddress)
					addresses.GET("/", handler.GetAddresses)
					addresses.PUT("/:id", handler.UpdateAddress)
					addresses.DELETE("/:id", handler.DeleteAddress)
				}

				// 用戶偏好
				preferences := user.Group("/preferences")
				{
					preferences.GET("/", handler.GetPreference)
					preferences.PUT("/", handler.UpdatePreference)
				}
			}
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
