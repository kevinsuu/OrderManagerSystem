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
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/handler"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加載配置
	cfg := config.LoadConfig()

	// 初始化日誌
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	// 連接數據庫
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化 Redis
	rdb := initRedis(cfg)

	// 初始化依賴
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, rdb)
	paymentHandler := handler.NewHandler(paymentService)

	// 設置 Gin 路由
	router := setupRouter(paymentHandler)

	// 創建 HTTP 服務器
	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}

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

// initRedis 初始化 Redis 客戶端
func initRedis(cfg *config.Config) repository.RedisRepository {
	return repository.NewRedisRepository(cfg.Redis)
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
		payments := api.Group("/payments")
		{
			// 支付管理
			payments.POST("/", h.CreatePayment)
			payments.GET("/", h.ListPayments)
			payments.GET("/:id", h.GetPayment)
			payments.GET("/order/:orderId", h.GetPaymentByOrderID)
			payments.GET("/user/:userId", h.GetUserPayments)

			// 支付操作
			payments.POST("/:id/process", h.ProcessPayment)
			payments.POST("/:id/cancel", h.CancelPayment)
			payments.POST("/refund", h.RefundPayment)
		}
	}

	return router
}
