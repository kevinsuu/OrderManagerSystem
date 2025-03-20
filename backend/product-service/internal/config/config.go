package config

import (
	"log"
	"os"
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
			Address: getEnv("SERVER_ADDRESS", ":8081"),
		},
		Firebase: FirebaseConfig{
			CredentialsFile: os.Getenv("FIREBASE_CREDENTIALS"),
			ProjectID:       os.Getenv("FIREBASE_PROJECT_ID"),
			DatabaseURL:     os.Getenv("FIREBASE_DATABASE_URL"),
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
