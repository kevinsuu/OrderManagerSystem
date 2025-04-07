package config

import (
	"log"
	"os"
)

// Config 應用配置
type Config struct {
	Server struct {
		Address string
	}
	Firebase struct {
		CredentialsFile string
		ProjectID       string
	}
	JWT struct {
		Secret string
	}
	ProductService struct {
		BaseURL string
	}
	OrderService struct {
		BaseURL string
	}
}

// ServerConfig 服務器配置
type ServerConfig struct {
	Address string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string
}

// ProductServiceConfig 產品服務配置
type ProductServiceConfig struct {
	BaseURL string
}

// OrderServiceConfig 訂單服務配置
type OrderServiceConfig struct {
	BaseURL string
}

// FirebaseConfig Firebase配置
type FirebaseConfig struct {
	CredentialsFile string
	ProjectID       string
}

// LoadConfig 加載配置
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", ":8082"),
		},
		Firebase: FirebaseConfig{
			CredentialsFile: os.Getenv("FIREBASE_CREDENTIALS"),
			ProjectID:       os.Getenv("FIREBASE_PROJECT_ID"),
		},
		ProductService: ProductServiceConfig{
			BaseURL: getEnv("PRODUCT_SERVICE_URL", "https://ordermanagersystem-product-service.onrender.com"),
		},
	}
}

// getEnv 獲取環境變量，如果不存在則返回默認值
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
