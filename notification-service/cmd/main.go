package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/handler"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/notification-service/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加載配置
	cfg := config.LoadConfig()

	// 初始化資料庫連接
	db := repository.NewPostgresDB(cfg.Database)
	// 初始化 Redis
	rdb := repository.NewRedisRepository(cfg.Redis)

	// 初始化存儲層
	notificationRepo := repository.NewNotificationRepository(db, rdb)

	// 初始化服務層
	notificationService := service.NewNotificationService(notificationRepo)

	// 初始化 HTTP 處理器
	notificationHandler := handler.NewHandler(notificationService)

	// 設置 Gin 路由
	router := setupRouter(notificationHandler)

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

// initDB 初始化數據庫連接
func initDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 設置連接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeMinutes) * time.Minute)

	return db, nil
}

// setupRouter 設置路由
func setupRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	// 中間件
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// 健康檢查
	router.GET("/health", h.HealthCheck)

	// API 路由組
	api := router.Group("/api/v1")
	{
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
