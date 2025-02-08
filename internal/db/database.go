package db

import (
	"fmt"
	"os"

	"memo-board/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	// .env 파일 로드
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using default environment variables")
	}

	// 개별 환경 변수에서 값 읽기
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	host := os.Getenv("MYSQL_HOST") // docker-compose에서는 "db", 로컬에서는 "localhost"
	port := os.Getenv("MYSQL_PORT") // 기본값 "3306" 사용

	if user == "" || password == "" || dbName == "" || host == "" || port == "" {
		return fmt.Errorf("MYSQL_USER, MYSQL_PASSWORD, MYSQL_DATABASE, MYSQL_HOST, MYSQL_PORT 환경 변수가 설정되지 않았습니다")
	}

	// DSN 구성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

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
