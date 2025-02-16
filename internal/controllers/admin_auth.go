package controllers

import (
	"memo-board/internal/db"
	"memo-board/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AdminLoginPage renders the admin login template.
func AdminLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_login.html", gin.H{})
}

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

// AdminApproveUser PUT /users/:id/approve
func AdminApproveUser(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminID := userIDVal.(uint)
	var admin models.User
	if err := db.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if admin.Role != models.ADMIN {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	idParam := c.Param("id")
	parsedID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, uint(parsedID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.IsApproved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already approved"})
		return
	}

	user.IsApproved = true
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User approved successfully"})
}
