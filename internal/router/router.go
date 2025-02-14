package router

import (
	"memo-board/internal/controllers"
	"memo-board/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// HTML 템플릿 로드
	r.LoadHTMLGlob("templates/*")

	// CORS 설정 (개발 시)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// 사용자 관련 라우트
	r.POST("/users/signup", controllers.UserSignup)
	r.POST("/users/login", controllers.UserLogin)

	// 게시글 관련 라우트
	r.GET("/posts", controllers.ListPosts)
	r.GET("/posts/:id", controllers.GetPost)

	// 인증 필요한 게시글 라우트
	auth := r.Group("/posts")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.POST("/", controllers.CreatePost)
		auth.PUT("/:id", controllers.UpdatePost)
		auth.DELETE("/:id", controllers.DeletePost)
	}

	// Admin 관련 라우트 그룹
	admin := r.Group("/admin")
	{
		admin.GET("/", controllers.AdminDashboard)
		admin.GET("/users", controllers.AdminUserList)
		admin.GET("/users/:id", controllers.AdminUserDetail)
		admin.POST("/users/:id", controllers.AdminUserUpdate)
	}

	return r
}
