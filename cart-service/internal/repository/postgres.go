package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresDB 創建並初始化 PostgreSQL 數據庫連接
func NewPostgresDB(cfg config.DatabaseConfig) *gorm.DB {
	// 首先連接到 postgres 默認數據庫
	defaultDSN := getDefaultDSN(cfg.DSN)
	defaultDB, err := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to default database: %v", err)
	}

	// 檢查數據庫是否存在
	dbName := "cart_service"
	var exists bool
	err = defaultDB.Raw("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = ?)", dbName).Scan(&exists).Error
	if err != nil {
		log.Fatalf("Failed to check database existence: %v", err)
	}

	// 只在數據庫不存在時創建
	if !exists {
		err = defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		log.Printf("Successfully created database: %s", dbName)
	}

	// 連接到應用數據庫
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 設置連接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMinutes) * time.Minute)

	// 自動遷移數據庫結構
	if err := db.AutoMigrate(&model.Cart{}, &model.CartItem{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Connected to PostgreSQL database and completed migrations")
	return db
}

// getDefaultDSN 獲取連接到默認數據庫的 DSN
func getDefaultDSN(originalDSN string) string {
	return "host=localhost user=postgres password=password dbname=postgres sslmode=disable"
}
