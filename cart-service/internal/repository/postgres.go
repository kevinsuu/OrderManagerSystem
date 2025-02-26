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
	if err := autoMigrate(db); err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	return db
}

// autoMigrate 自動遷移數據庫結構
func autoMigrate(db *gorm.DB) error {
	// 添加您需要遷移的模型
	if err := db.AutoMigrate(&model.Cart{}, &model.CartItem{}); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}
	return nil
}
