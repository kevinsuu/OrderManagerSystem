package repository

import (
	"fmt"
	"log"

	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/config"
	"github.com/kevinsuu/OrderManagerSystem/auth-service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresDB 創建新的 PostgreSQL 數據庫連接
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
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetimeMinutes)

	// 自動遷移數據庫結構
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL database")
	return db
}
