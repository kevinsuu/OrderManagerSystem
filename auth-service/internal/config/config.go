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
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret        string
	ExpiryMinutes int
}

// LoadConfig 加載配置
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Address: ":8083",
		},
		Firebase: FirebaseConfig{
			CredentialsFile: os.Getenv("FIREBASE_CREDENTIALS"),
			ProjectID:       os.Getenv("FIREBASE_PROJECT_ID"),
		},
		JWT: JWTConfig{
			Secret:        os.Getenv("JWT_SECRET"),
			ExpiryMinutes: getEnvAsInt("JWT_TOKEN_EXPIRY_MINUTES", 60),
		},
	}
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
