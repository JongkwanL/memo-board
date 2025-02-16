package middleware

import (
	"net/http"

	"memo-board/internal/db"
	"memo-board/internal/models"

	"github.com/gin-gonic/gin"
)

// AdminJWTAuthMiddleware : JWTAuthMiddleware 이후에 실행되어, 컨텍스트에 저장된 user_id를 이용해
func AdminJWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID := userIDVal.(uint)
		var user models.User
		if err := db.DB.First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if user.Role != models.ADMIN {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin only"})
			return
		}

		c.Next()
	}
}
