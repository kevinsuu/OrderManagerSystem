package config

import (
	"log"
	"os"
	"strconv"
)

// Config 應用配置
type Config struct {
	Server         ServerConfig
	Database       DatabaseConfig
	Redis          RedisConfig
	JWT            JWTConfig
	ProductService ProductServiceConfig
	OrderService   OrderServiceConfig
}

// ServerConfig 服務器配置
type ServerConfig struct {
	Address string
}

// DatabaseConfig 數據庫配置
type DatabaseConfig struct {
	DSN                    string
	MaxIdleConns           int
	MaxOpenConns           int
	ConnMaxLifetimeMinutes int
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
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

// LoadConfig 加載配置
func LoadConfig() *Config {
	return &Config{
		Server:         loadServerConfig(),
		Database:       loadDatabaseConfig(),
		Redis:          loadRedisConfig(),
		JWT:            loadJWTConfig(),
		ProductService: loadProductServiceConfig(),
		OrderService:   loadOrderServiceConfig(),
	}
}

// loadServerConfig 加載服務器配置
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Address: getEnv("SERVER_ADDRESS", ":8086"),
	}
}

// loadDatabaseConfig 加載數據庫配置
func loadDatabaseConfig() DatabaseConfig {
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "100"))
	connMaxLifetime, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME_MINUTES", "60"))

	return DatabaseConfig{
		DSN:                    getDatabaseDSN(),
		MaxIdleConns:           maxIdleConns,
		MaxOpenConns:           maxOpenConns,
		ConnMaxLifetimeMinutes: connMaxLifetime,
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

// loadJWTConfig 加載JWT配置
func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret: getEnv("JWT_SECRET", "your-secret-key"),
	}
}

// loadOrderServiceConfig 加載訂單服務配置
func loadOrderServiceConfig() OrderServiceConfig {
	return OrderServiceConfig{
		BaseURL: getEnv("ORDER_SERVICE_BASE_URL", "http://localhost:8082"),
	}
}

// loadProductServiceConfig 加載產品服務配置
func loadProductServiceConfig() ProductServiceConfig {
	return ProductServiceConfig{
		BaseURL: getEnv("PRODUCT_SERVICE_BASE_URL", "http://localhost:8081"),
	}
}

// getDatabaseDSN 獲取數據庫連接字符串
func getDatabaseDSN() string {
	// 如果提供了完整的 DSN，直接使用
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		return dsn
	}

	// 否則從各個組件構建 DSN
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "cart_service")
	sslmode := getEnv("DB_SSLMODE", "disable")

	return "host=" + host + " port=" + port + " user=" + user + " password=" + password +
		" dbname=" + dbname + " sslmode=" + sslmode
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
