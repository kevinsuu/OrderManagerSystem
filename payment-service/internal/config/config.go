package config

import (
	"fmt"
	"log"
	"os"
)

// Config 應用配置
type Config struct {
	Server   ServerConfig
	Firebase FirebaseConfig
	JWT      JWTConfig
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string
}

// ServerConfig 服務器配置
type ServerConfig struct {
	Address string
}

// FirebaseConfig Firebase 配置
type FirebaseConfig struct {
	CredentialsFile string
	ProjectID       string
	DatabaseURL     string
}

// LoadConfig 加載配置
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", ":8084"),
		},
		Firebase: FirebaseConfig{
			CredentialsFile: os.Getenv("FIREBASE_CREDENTIALS"),
			ProjectID:       os.Getenv("FIREBASE_PROJECT_ID"),
		},
	}
}

// validate 驗證配置
func (c *Config) validate() error {
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET environment variable must be set")
	}
	log.Printf("JWT Secret length: %d", len(c.JWT.Secret))
	return nil
}

// loadJWTConfig 加載 JWT 配置
func loadJWTConfig() JWTConfig {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Printf("Warning: JWT_SECRET environment variable is not set")
	}
	return JWTConfig{
		Secret: secret,
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
