package controllers

import (
	"memo-board/internal/db"
	"memo-board/internal/middleware"
	"memo-board/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// POST /users/signup
func UserSignup(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email"    binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 비밀번호 해싱
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)

	user := models.User{
		Username: req.Username,
		Password: string(hashed),
		Email:    req.Email,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "signup success"})
}

// POST /users/login
func UserLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 비번 검증
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// JWT 발급
	token, err := middleware.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// PUT /users/:id/approve
func AdminApproveUser(c *gin.Context) {
	// 미들웨어를 통해 JWT에서 추출된 사용자 정보를 컨텍스트에 "user" 키로 저장했다고 가정합니다.
	// 이를 이용하여 호출한 사용자가 admin인지 확인합니다.
	adminUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userData, ok := adminUser.(models.User)
	if !ok || userData.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	// 승인할 사용자 ID를 URL 파라미터로부터 추출
	id := c.Param("id")
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 이미 승인된 경우 확인 (선택 사항)
	if user.IsApproved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already approved"})
		return
	}

	// 승인 처리: IsApproved를 true로 업데이트
	user.IsApproved = true
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User approved successfully"})
}
