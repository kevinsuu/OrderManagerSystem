package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware 日誌中間件
func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 開始時間
		startTime := time.Now()

		// 處理請求
		c.Next()

		// 結束時間
		endTime := time.Now()

		// 執行時間
		latencyTime := endTime.Sub(startTime)

		// 請求方法
		reqMethod := c.Request.Method

		// 請求路由
		reqUri := c.Request.RequestURI

		// 狀態碼
		statusCode := c.Writer.Status()

		// 請求IP
		clientIP := c.ClientIP()

		// 日誌格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info("HTTP Request")
	}
}
