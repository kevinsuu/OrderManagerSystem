package config

import (
	"log"
	"os"
)

// Config 應用配置
type Config struct {
	Server   ServerConfig
	Services ServicesConfig
}

// ServerConfig 服務器配置
type ServerConfig struct {
	Port string
}

// ServicesConfig 微服務配置
type ServicesConfig struct {
	AuthService         string
	OrderService       string
	ProductService     string
	PaymentService     string
	NotificationService string
}

// LoadConfig 加載配置
func LoadConfig() *Config {
	return &Config{
		Server:   loadServerConfig(),
		Services: loadServicesConfig(),
	}
}

// loadServerConfig 加載服務器配置
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port: getEnv("SERVER_PORT", "8080"),
	}
}

// loadServicesConfig 加載微服務配置
func loadServicesConfig() ServicesConfig {
	return ServicesConfig{
		AuthService:         getEnv("AUTH_SERVICE_URL", "http://auth-service:8081"),
		OrderService:       getEnv("ORDER_SERVICE_URL", "http://order-service:8082"),
		ProductService:     getEnv("PRODUCT_SERVICE_URL", "http://product-service:8083"),
		PaymentService:     getEnv("PAYMENT_SERVICE_URL", "http://payment-service:8084"),
		NotificationService: getEnv("NOTIFICATION_SERVICE_URL", "http://notification-service:8085"),
	}
}

// getEnv 獲取環境變量
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue == "" {
			log.Printf("Warning: Environment variable %s not set and no default value provided", key)
		}
		return defaultValue
	}
	return value
}
