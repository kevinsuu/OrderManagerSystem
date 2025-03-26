package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/handler"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/infrastructure/firebase"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/middleware"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/service"
)

func main() {
	// 加載配置
	cfg := config.LoadConfig()

	// 初始化日誌
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	// 初始化 Firebase
	ctx := context.Background()
	fb, err := firebase.InitFirebase(ctx, cfg.Firebase.CredentialsFile, cfg.Firebase.ProjectID)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// 初始化存儲層
	notificationRepo := repository.NewNotificationRepository(fb.Database)

	// 初始化服務層
	notificationService := service.NewNotificationService(notificationRepo)

	// 初始化 HTTP 處理器
	notificationHandler := handler.NewHandler(notificationService)

	// 設置 Gin 路由
	router := setupRouter(notificationHandler, cfg.JWT.Secret)

	// 創建 HTTP 服務器
	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}

	// 啟動通知處理器
	go startNotificationProcessor(notificationService)

	// 在後台運行服務器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 設置關閉超時
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 關閉 HTTP 服務器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// setupRouter 設置路由
func setupRouter(h *handler.Handler, jwtSecret string) *gin.Engine {
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

	// 健康檢查
	router.GET("/health", h.HealthCheck)

	// API 路由組
	api := router.Group("/api/v1")
	{
		// 添加認證中間件
		api.Use(middleware.AuthMiddleware(jwtSecret))

		notifications := api.Group("/notifications")
		{
			// 通知管理
			notifications.POST("/", h.CreateNotification)
			notifications.POST("/template", h.CreateNotificationFromTemplate)
			notifications.GET("/", h.ListNotifications)
			notifications.GET("/:id", h.GetNotification)
			notifications.GET("/user/:userId", h.GetUserNotifications)
		}

		templates := api.Group("/templates")
		{
			// 模板管理
			templates.POST("/", h.CreateTemplate)
			templates.GET("/", h.ListTemplates)
			templates.GET("/:id", h.GetTemplate)
			templates.PUT("/:id", h.UpdateTemplate)
		}
	}

	return router
}

// startNotificationProcessor 啟動通知處理器
func startNotificationProcessor(notificationService service.NotificationService) {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		if err := notificationService.ProcessPendingNotifications(context.Background()); err != nil {
			log.Printf("Failed to process pending notifications: %v", err)
		}
	}
}
