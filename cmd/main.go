package main

import (
	"log"
	"memo-board/internal/db"
	"memo-board/internal/router"
	"os"
)

func main() {
	// (선택) .env 파일 로드 (로컬 개발 시)
	// err := godotenv.Load()
	// if err != nil {
	//     log.Println("Warning: No .env file found")
	// }

	// DB 초기화
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to init DB: %v\n", err)
	}

	// 라우터 설정
	r := router.SetupRouter()

	// 포트 설정
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	r.Run(":" + port)
}
