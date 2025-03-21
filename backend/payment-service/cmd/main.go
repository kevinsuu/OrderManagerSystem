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
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/handler"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/infrastructure/firebase"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/middleware"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/repository"
	"github.com/kevinsuu/OrderManagerSystem/payment-service/internal/service"
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

	// 初始化依賴
	paymentRepo := repository.NewPaymentRepository(fb.Database)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentHandler := handler.NewHandler(paymentService)

	// 設置路由
	router := gin.Default()

	// 中間件
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// CORS 中間件配置
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// 健康檢查路由（不需要認證）
	router.GET("/health", paymentHandler.HealthCheck)

	// API 路由組
	api := router.Group("/api/v1")
	{
		// 添加認證中間件
		api.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

		payments := api.Group("/payments")
		{
			// 支付管理（按具體到通用的順序排列）
			payments.GET("/order/:orderId", paymentHandler.GetPaymentByOrderID) // 最具體的路由放在前面
			payments.GET("/user/:userId", paymentHandler.GetUserPayments)
			payments.POST("/refund", paymentHandler.RefundPayment)
			payments.POST("/:id/process", paymentHandler.ProcessPayment)
			payments.POST("/:id/cancel", paymentHandler.CancelPayment)
			payments.POST("/", paymentHandler.CreatePayment)
			payments.GET("/", paymentHandler.ListPayments)
			payments.GET("/:id", paymentHandler.GetPayment) // 最通用的路由放在最後
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
