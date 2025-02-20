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
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/middleware"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/service"
)

func main() {
	// 加載配置
	cfg := config.LoadConfig()

	// 初始化日誌
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	// 連接數據庫
	db := repository.NewPostgresDB(cfg.Database)

	// 初始化 Redis
	rdb := initRedis(cfg)

	// 初始化依賴
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, rdb)
	paymentHandler := handler.NewHandler(paymentService)

	// 設置路由
	router := gin.Default()

	// 健康檢查路由（不需要認證）
	router.GET("/health", paymentHandler.HealthCheck)

	// API 路由組
	api := router.Group("/api/v1")
	{
		// 添加認證中間件
		api.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

		payments := api.Group("/payments")
		{
			// 支付管理
			payments.POST("/", paymentHandler.CreatePayment)
			payments.GET("/", paymentHandler.ListPayments)
			payments.GET("/:id", paymentHandler.GetPayment)
			payments.GET("/order/:orderId", paymentHandler.GetPaymentByOrderID)
			payments.GET("/user/:userId", paymentHandler.GetUserPayments)

			// 支付操作
			payments.POST("/:id/process", paymentHandler.ProcessPayment)
			payments.POST("/:id/cancel", paymentHandler.CancelPayment)
			payments.POST("/refund", paymentHandler.RefundPayment)
		}
	}

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

// initRedis 初始化 Redis 客戶端
func initRedis(cfg *config.Config) repository.RedisRepository {
	return repository.NewRedisRepository(cfg.Redis)
}
