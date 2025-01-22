package db

import (
	"fmt"
	"os"

	"memo-board/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// 예: MySQL
		dsn = "root:password@tcp(localhost:3306)/board?charset=utf8mb4&parseTime=True&loc=Local"
		// 또는 PostgreSQL
		// dsn = "postgres://user:pass@localhost:5432/board?sslmode=disable"
	}

	// 드라이버에 따라 선택
	// dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect DB: %w", err)
	}

	// 자동 마이그레이션
	if err := dbConn.AutoMigrate(&models.User{}, &models.Post{}); err != nil {
		return fmt.Errorf("auto migrate error: %w", err)
	}

	DB = dbConn
	return nil
}
