package auth

import (
	"log" // <--- [MỚI] Thêm thư viện này để in lỗi ra màn hình
	"mangahub/internal/user"
	"mangahub/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	UserRepo *user.Repository
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Register(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Tạo User ID ngẫu nhiên
	newUser := &user.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Password: req.Password,
	}

	// [QUAN TRỌNG] Đoạn này đã được sửa để in lỗi ra
	if err := h.UserRepo.CreateUser(newUser); err != nil {
		// 1. In lỗi chi tiết ra Terminal của Server để bạn đọc
		log.Printf("❌ LỖI DATABASE: %v", err)

		// 2. Trả về lỗi chi tiết cho Client (curl) xem luôn
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Could not create user",
			"detail": err.Error(), // <--- Dòng này sẽ cho bạn biết chính xác lý do
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 1. Tìm user trong DB
	u, err := h.UserRepo.GetUserByUsername(req.Username)
	if err != nil {
		// In lỗi ra terminal nếu login fail (để debug)
		log.Printf("⚠️ Login failed for user %s: %v", req.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 2. So khớp mật khẩu
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 3. Tạo JWT Token
	token, err := utils.GenerateToken(u.ID, u.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	log.Println("GENERATED TOKEN:", token)

	// 4. Trả về token cho client
	c.JSON(http.StatusOK, gin.H{"token": token})
}