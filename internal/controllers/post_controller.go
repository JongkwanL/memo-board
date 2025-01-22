package controllers

import (
	"memo-board/internal/db"
	"memo-board/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST /posts (JWT 필요)
func CreatePost(c *gin.Context) {
	userIDany, _ := c.Get("user_id") // 미들웨어서 set
	userID := userIDany.(uint)

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: userID,
	}
	if err := db.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

// GET /posts/:id
func GetPost(c *gin.Context) {
	idParam := c.Param("id")
	var post models.Post

	if err := db.DB.Preload("Author").
		Where("id = ?", idParam).
		First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// GET /posts (목록)
func ListPosts(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var posts []models.Post
	db.DB.Preload("Author").Order("id DESC").
		Limit(limit).Offset(offset).Find(&posts)

	c.JSON(http.StatusOK, posts)
}

// PUT /posts/:id (JWT 필요, 작성자만)
func UpdatePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	idParam := c.Param("id")

	var post models.Post
	if err := db.DB.Where("id = ?", idParam).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	// 작성자 본인만 수정 허용
	if post.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not owner"})
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.Title = req.Title
	post.Content = req.Content
	db.DB.Save(&post)

	c.JSON(http.StatusOK, post)
}

// DELETE /posts/:id (JWT 필요, 작성자만)
func DeletePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	idParam := c.Param("id")

	var post models.Post
	if err := db.DB.Where("id = ?", idParam).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	if post.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not owner"})
		return
	}

	db.DB.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
