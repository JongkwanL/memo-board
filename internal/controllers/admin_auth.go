package controllers

import (
	"memo-board/internal/db"
	"memo-board/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AdminLogin handles POST /admin/login.
func AdminLogin(c *gin.Context) {
	adminID := c.PostForm("admin_id")
	adminPW := c.PostForm("admin_pw")

	var admin models.User
	if err := db.DB.Where("username = ?", adminID).First(&admin).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "admin_login.html", gin.H{
			"Error": "아이디 또는 비밀번호가 잘못되었습니다.",
		})
		return
	}

	if admin.Role != models.ADMIN {
		c.HTML(http.StatusUnauthorized, "admin_login.html", gin.H{
			"Error": "관리자 계정이 아닙니다.",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(adminPW)); err != nil {
		c.HTML(http.StatusUnauthorized, "admin_login.html", gin.H{
			"Error": "아이디 또는 비밀번호가 잘못되었습니다.",
		})
		return
	}

	// 로그인 성공: 필요한 경우 세션이나 토큰을 생성한 후 관리자 대시보드로 리다이렉트합니다.
	c.Redirect(http.StatusFound, "/admin/dashboard")
}
