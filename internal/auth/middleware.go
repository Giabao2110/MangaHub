package auth

import (
	"mangahub/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware là "bác bảo vệ" kiểm tra Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy token từ header "Authorization"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 2. Cắt bỏ chữ "Bearer " (nếu có) để lấy token sạch
		// Định dạng chuẩn là: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Kiểm tra Token có hợp lệ không
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 4. Nếu ngon lành, lưu thông tin user vào context để các hàm sau dùng
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next() // Cho phép đi tiếp vào trong
	}
}
