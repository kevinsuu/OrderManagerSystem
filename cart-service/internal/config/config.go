package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server         ServerConfig
	Redis          RedisConfig
	Database       DatabaseConfig
	JWT            JWTConfig
	ProductService ProductServiceConfig
}

type ServerConfig struct {
	Address string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type DatabaseConfig struct {
	DSN                    string
	MaxIdleConns           int
	MaxOpenConns           int
	ConnMaxLifetimeMinutes int
}

type JWTConfig struct {
	Secret string
}

type ProductServiceConfig struct {
	BaseURL string
}

// LoadConfig 加載配置
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", ":8082"),
		},
		Redis:    loadRedisConfig(),
		Database: loadDatabaseConfig(),
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
		},
		ProductService: ProductServiceConfig{
			BaseURL: getEnv("PRODUCT_SERVICE_BASE_URL", "http://localhost:8080"),
		},
	}
}

// loadRedisConfig 加載Redis配置
func loadRedisConfig() RedisConfig {
	db, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		log.Printf("Warning: Invalid REDIS_DB value, using default (0)")
		db = 0
	}

	return RedisConfig{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       db,
	}
}

// loadDatabaseConfig 加載數據庫配置
func loadDatabaseConfig() DatabaseConfig {
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "100"))
	connMaxLifetimeMinutes, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME_MINUTES", "60"))

	return DatabaseConfig{
		DSN:                    getEnv("DB_DSN", "host=localhost user=postgres password=postgres dbname=cart_service port=5432 sslmode=disable"),
		MaxIdleConns:           maxIdleConns,
		MaxOpenConns:           maxOpenConns,
		ConnMaxLifetimeMinutes: connMaxLifetimeMinutes,
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
