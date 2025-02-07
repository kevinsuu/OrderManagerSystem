package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ServiceConfig 服務配置
type ServiceConfig struct {
	AuthService        string
	OrderService      string
	ProductService    string
	PaymentService    string
	NotificationService string
}

// ProxyHandler 代理處理器
type ProxyHandler struct {
	config ServiceConfig
}

// NewProxyHandler 創建代理處理器
func NewProxyHandler(config ServiceConfig) *ProxyHandler {
	return &ProxyHandler{
		config: config,
	}
}

// ReverseProxy 反向代理處理
func (h *ProxyHandler) ReverseProxy(targetURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse upstream url"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Param("proxyPath")
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// HealthCheck 健康檢查
func (h *ProxyHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"services": map[string]string{
			"auth":         h.config.AuthService,
			"order":        h.config.OrderService,
			"product":      h.config.ProductService,
			"payment":      h.config.PaymentService,
			"notification": h.config.NotificationService,
		},
	})
}
