package config

import (
	"log"
	"os"
	"strconv"
)

// Config 應用配置
type Config struct {
	Server   ServerConfig
	Firebase FirebaseConfig
	JWT      JWTConfig
}

// ServerConfig 服務器配置
type ServerConfig struct {
	Address string
}

// FirebaseConfig Firebase配置
type FirebaseConfig struct {
	CredentialsFile string
	ProjectID       string
	DatabaseURL     string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string
}

// LoadConfig 加載配置
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", ":8085"),
		},
		Firebase: FirebaseConfig{
			CredentialsFile: os.Getenv("FIREBASE_CREDENTIALS"),
			ProjectID:       os.Getenv("FIREBASE_PROJECT_ID"),
			DatabaseURL:     os.Getenv("FIREBASE_DATABASE_URL"),
		},
		JWT: loadJWTConfig(),
	}
}

// loadJWTConfig 加載 JWT 配置
func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret: os.Getenv("JWT_SECRET"),
	}
}

// validate 驗證配置
func (c *Config) validate() {
	// 驗證 Firebase 配置
	if c.Firebase.CredentialsFile == "" {
		log.Fatal("FIREBASE_CREDENTIALS environment variable is required")
	}
	if c.Firebase.ProjectID == "" {
		log.Fatal("FIREBASE_PROJECT_ID environment variable is required")
	}
	if c.Firebase.DatabaseURL == "" {
		log.Fatal("FIREBASE_DATABASE_URL environment variable is required")
	}

	// 驗證 JWT 配置
	if c.JWT.Secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// 記錄配置信息
	log.Printf("Server will run on: %s", c.Server.Address)
	log.Printf("Firebase Project ID: %s", c.Firebase.ProjectID)
	log.Printf("Using JWT secret with length: %d", len(c.JWT.Secret))
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

// getEnvAsInt 獲取環境變量，如果不存在則返回默認值
func getEnvAsInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}
