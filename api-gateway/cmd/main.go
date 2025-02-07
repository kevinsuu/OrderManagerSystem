package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 健康檢查端點
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
        })
    })

    // 從環境變數獲取服務 URL
    authServiceURL := os.Getenv("AUTH_SERVICE_URL")
    orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
    productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
    paymentServiceURL := os.Getenv("PAYMENT_SERVICE_URL")
    notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")

    // 設置路由組
    auth := r.Group("/auth")
    {
        auth.Any("/*path", proxyHandler(authServiceURL))
    }

    orders := r.Group("/orders")
    {
        orders.Any("/*path", proxyHandler(orderServiceURL))
    }

    products := r.Group("/products")
    {
        products.Any("/*path", proxyHandler(productServiceURL))
    }

    payments := r.Group("/payments")
    {
        payments.Any("/*path", proxyHandler(paymentServiceURL))
    }

    notifications := r.Group("/notifications")
    {
        notifications.Any("/*path", proxyHandler(notificationServiceURL))
    }

    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func proxyHandler(targetURL string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 在這裡實現代理邏輯
        c.JSON(http.StatusOK, gin.H{
            "message": "API Gateway is working",
            "target_url": targetURL,
            "path": c.Param("path"),
        })
    }
}
