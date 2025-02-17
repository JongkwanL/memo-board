package main

import (
	"log"
	"memo-board/internal/db"
	"memo-board/internal/router"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// .env 파일 로드 (로컬 개발 시)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// DB 초기화
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to init DB: %v\n", err)
	}

	// 라우터 설정 (내부 router 패키지에서 admin 관련 엔드포인트 포함)
	r := router.SetupRouter()

	// 포트 설정
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	r.Run(":" + port)
}
