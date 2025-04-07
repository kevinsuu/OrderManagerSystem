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

	// 使用 gin.New() 而不是 gin.Default()
	router := gin.New()

	// CORS 中間件配置
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Logger 和 Recovery 中間件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	// API 路由組
	api := router.Group("/api/v1")
	{
		// 公開路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
			auth.POST("/refresh", handler.RefreshToken)
			auth.POST("/forgot-password", handler.ForgotPassword)
			auth.POST("/reset-password", handler.ResetPassword)
		}

		// 需要認證的路由
		secured := api.Group("/user")
		secured.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			//取得使用者資訊
			secured.GET("/", handler.GetUser)
			// 用戶偏好設置
			secured.GET("/preferences", handler.GetPreference)
			secured.PUT("/preferences", handler.UpdatePreference)

			// 地址管理
			secured.GET("/addresses", handler.GetAddresses)
			secured.POST("/addresses", handler.CreateAddress)
			secured.PUT("/addresses/:id", handler.UpdateAddress)
			secured.DELETE("/addresses/:id", handler.DeleteAddress)
			secured.PUT("/addresses/:id/default", handler.SetDefaultAddress)
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
