package router

import (
	"memo-board/internal/controllers"
	"memo-board/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.POST("/users/signup", controllers.UserSignup)
	r.POST("/users/login", controllers.UserLogin)

	r.GET("/posts", controllers.ListPosts)
	r.GET("/posts/:id", controllers.GetPost)

	auth := r.Group("/posts")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.POST("/", controllers.CreatePost)
		auth.PUT("/:id", controllers.UpdatePost)
		auth.DELETE("/:id", controllers.DeletePost)
	}

	adminPublic := r.Group("/admin")
	{
		adminPublic.GET("/login", controllers.AdminLoginPage)
		adminPublic.POST("/login", controllers.AdminLogin)
		adminPublic.GET("/dashboard", controllers.AdminDashboard)
		adminPublic.GET("/users", controllers.AdminUserList)

	}

	adminProtected := r.Group("/admin")
	adminProtected.Use(middleware.JWTAuthMiddleware(), middleware.AdminJWTAuthMiddleware())
	{
		adminProtected.GET("/dashboard-data", controllers.AdminDashboardData)
		adminProtected.GET("/users/:id", controllers.AdminUserDetail)
		adminProtected.POST("/users/:id", controllers.AdminUserUpdate)
		adminProtected.PUT("/users/:id/approve", controllers.AdminApproveUser)
		adminProtected.GET("/users-data", controllers.AdminUsersData)
	}

	return r
}
