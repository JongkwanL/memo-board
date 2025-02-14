package controllers

import (
	"github.com/gin-gonic/gin"
	"memo-board/internal/db"
	"memo-board/internal/models"
	"net/http"
)

// AdminUserList handles GET /admin/users
func AdminUserList(c *gin.Context) {
	usernameFilter := c.Query("username")
	roleFilter := c.Query("role")
	isApprovedFilter := c.Query("isapproved")

	query := db.DB.Model(&models.User{})
	if usernameFilter != "" {
		query = query.Where("username LIKE ?", "%"+usernameFilter+"%")
	}
	if roleFilter != "" {
		query = query.Where("role = ?", roleFilter)
	}
	if isApprovedFilter != "" {
		if isApprovedFilter == "true" {
			query = query.Where("is_approved = ?", true)
		} else if isApprovedFilter == "false" {
			query = query.Where("is_approved = ?", false)
		}
	}

	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error: %v", err)
		return
	}

	c.HTML(http.StatusOK, "admin_users.html", gin.H{
		"Users":            users,
		"UsernameFilter":   usernameFilter,
		"RoleFilter":       roleFilter,
		"IsApprovedFilter": isApprovedFilter,
	})
}

// AdminUserDetail handles GET /admin/users/:id
func AdminUserDetail(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}
	c.HTML(http.StatusOK, "admin_user_detail.html", gin.H{
		"User": user,
	})
}

// AdminUserUpdate handles POST /admin/users/:id
func AdminUserUpdate(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	isApproved := c.PostForm("isapproved")
	user.IsApproved = (isApproved == "on")

	if err := db.DB.Save(&user).Error; err != nil {
		c.String(http.StatusInternalServerError, "Update failed: %v", err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/users")
}

// AdminDashboard handles GET /admin
func AdminDashboard(c *gin.Context) {
	// 예시: /admin 접속 시 유저 리스트 페이지로 리다이렉트
	c.Redirect(http.StatusFound, "/admin/users")
}
