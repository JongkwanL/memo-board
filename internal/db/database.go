package db

import (
	"errors"
	"fmt"
	"os"

	"memo-board/internal/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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

	// 관리 계정 생성 (환경 변수에서 ADMIN_ID, ADMIN_PW 읽기)
	adminID := os.Getenv("ADMIN_ID")
	adminPW := os.Getenv("ADMIN_PW")
	if adminID != "" && adminPW != "" {
		var admin models.User
		// admin 계정이 존재하는지 확인
		if err := dbConn.Where("username = ?", adminID).First(&admin).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 비밀번호 해싱
				hashed, err := bcrypt.GenerateFromPassword([]byte(adminPW), 14)
				if err != nil {
					return fmt.Errorf("failed to hash admin password: %w", err)
				}
				// admin 계정 생성, 이메일은 기본값으로 지정
				admin = models.User{
					Username:   adminID,
					Password:   string(hashed),
					Email:      adminID + "@admin.com",
					Role:       "admin",
					IsApproved: true,
				}
				if err := dbConn.Create(&admin).Error; err != nil {
					return fmt.Errorf("failed to create admin user: %w", err)
				}
				fmt.Println("Admin account created successfully")
			} else {
				return fmt.Errorf("failed to query admin user: %w", err)
			}
		}
	}

	DB = dbConn
	return nil
}
