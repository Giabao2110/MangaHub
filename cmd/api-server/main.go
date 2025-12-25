package main

import (
	"fmt"
	"log"
	"os" // Dùng os thay cho ioutil để hết bị gạch chéo
	"sort"
	"strings"
	"time"

	"mangahub/pkg/database"
	"mangahub/internal/chat" 
	"mangahub/pkg/utils"




	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	//"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.Und)

func main() {
	// =====================
	// INIT DATABASE
	// =====================
	db := database.InitDB("./mangahub.db")
	defer db.Close()
	chat.InitChatDB(db)
	log.Println("✅ SQLite connected")

	// =====================
	// INIT GIN
	// =====================
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// =====================
	// STATIC FILES
	// =====================
	r.Static("/static", "./web/static")
	r.Static("/data", "./data")

	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	r.GET("/login", func(c *gin.Context) {
		c.File("./web/login.html")
	})

	r.GET("/register", func(c *gin.Context) {
		c.File("./web/register.html")
	})

	// =====================
	// AUTH APIs (POST)
	// =====================

	// REGISTER - Đăng ký tài khoản
	r.POST("/api/register", func(c *gin.Context) {
		var u struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
			return
		}

		// Hash mật khẩu an toàn
		hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

		_, err := db.Exec(
			"INSERT INTO users (id, username, password_hash) VALUES (?, ?, ?)",
			u.Username,
			u.Username,
			string(hash),
		)

		if err != nil {
			c.JSON(400, gin.H{"error": "Username đã tồn tại"})
			return
		}

		c.JSON(200, gin.H{"message": "Đăng ký thành công"})
	})

	// LOGIN - Đăng nhập & Tạo Token
	r.POST("/api/login", func(c *gin.Context) {
		var u struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
			return
		}


		var id, hash, role string
		err := db.QueryRow(
			"SELECT id, password_hash, role FROM users WHERE username=?",
			u.Username,
		).Scan(&id, &hash, &role)

		if err != nil {
			c.JSON(401, gin.H{"error": "Tài khoản không tồn tại"})
			return
		}

		// So sánh mật khẩu
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password)); err != nil {
			c.JSON(401, gin.H{"error": "Sai mật khẩu"})
			return
		}

		// Tạo JWT Token
		tokenString, err := utils.GenerateToken(id, u.Username)
		if err != nil {
			c.JSON(500, gin.H{"error": "Không thể tạo token"})
			return
		}

		c.JSON(200, gin.H{
			"token":    tokenString,
			"username": u.Username,
			"role":     role,
		})
	})

	// =====================
	// MANGA APIs
	// =====================

	// LIST MANGAS - Lấy danh sách truyện
	r.GET("/api/mangas", func(c *gin.Context) {
		mangas := []gin.H{}

		files, err := os.ReadDir("./data")
		if err != nil {
			c.JSON(200, mangas)
			return
		}

		for _, f := range files {
			if f.IsDir() {
				slug := f.Name()
				title := titleCaser.String(strings.ReplaceAll(slug, "-", " "))

				// LOGIC QUÉT ẢNH THÔNG MINH:
				// Tìm file ảnh đầu tiên trong thư mục để làm ảnh bìa
				coverImg := ""
				subFiles, _ := os.ReadDir("./data/" + slug)
				for _, sf := range subFiles {
					name := strings.ToLower(sf.Name())
					// Kiểm tra xem có phải file ảnh không
					if !sf.IsDir() && (strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpeg")) {
						coverImg = sf.Name()
						break // Lấy file ảnh đầu tiên tìm thấy làm Cover
					}
				}

				// Nếu không thấy ảnh nào, dùng ảnh mặc định để tránh lỗi giao diện
				finalCover := fmt.Sprintf("/data/%s/%s", slug, coverImg)
				if coverImg == "" {
					finalCover = "https://via.placeholder.com/200x300?text=No+Cover"
				}

				mangas = append(mangas, gin.H{
					"slug":     slug,
					"title":    title,
					"cover":    finalCover,
					"category": "Manga",
				})
			}
		}

		c.JSON(200, mangas)
	})

	// READ CHAPTER - Lấy danh sách ảnh trong Chapter
	r.GET("/api/read/:slug/:chap", func(c *gin.Context) {
		dir := fmt.Sprintf("./data/%s/%s", c.Param("slug"), c.Param("chap"))

		files, err := os.ReadDir(dir)
		if err != nil {
			c.JSON(404, gin.H{"error": "Không tìm thấy chapter"})
			return
		}

		images := []string{}
		for _, f := range files {
			name := strings.ToLower(f.Name())
			if !f.IsDir() && (strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpeg")) {
				images = append(images, fmt.Sprintf("/data/%s/%s/%s", c.Param("slug"), c.Param("chap"), f.Name()))
			}
		}

		sort.Strings(images)
		c.JSON(200, images)
	})

	// 1. Lấy danh sách Chapter
    r.GET("/api/manga/:slug", func(c *gin.Context) {
        slug := c.Param("slug")
        path := fmt.Sprintf("./data/%s", slug)
        files, err := os.ReadDir(path)
        if err != nil {
            c.JSON(404, gin.H{"error": "Không tìm thấy truyện"})
            return
        }
        chapters := []string{}
        for _, f := range files {
            if f.IsDir() { chapters = append(chapters, f.Name()) }
        }
        sort.Strings(chapters)
        c.JSON(200, gin.H{"slug": slug, "title": titleCaser.String(strings.ReplaceAll(slug, "-", " ")), "chapters": chapters})
    })

    // 2. Thêm vào Wishlist (Yêu thích)
    r.POST("/api/wishlist", func(c *gin.Context) {
        var req struct {
            Username string `json:"username"`
            Slug     string `json:"slug"`
        }
        if err := c.ShouldBindJSON(&req); err != nil { return }
        _, err := db.Exec("INSERT INTO wishlist (username, manga_slug) VALUES (?, ?)", req.Username, req.Slug)
        if err != nil {
            c.JSON(400, gin.H{"error": "Đã có trong danh sách yêu thích"})
            return
        }
        c.JSON(200, gin.H{"message": "Đã thêm vào Wishlist"})
    })

    // 3. Gửi tin nhắn cho Admin
    r.POST("/api/messages", func(c *gin.Context) {
        var msg struct {
            Username string `json:"username"`
            Content  string `json:"content"`
        }
        c.ShouldBindJSON(&msg)
        _, err := db.Exec("INSERT INTO messages (username, content, created_at) VALUES (?, ?, ?)", 
            msg.Username, msg.Content, time.Now())
        if err != nil {
            c.JSON(500, gin.H{"error": "Lỗi gửi tin nhắn"})
            return
        }
        c.JSON(200, gin.H{"message": "Đã gửi thành công"})
    })

	// API lấy danh sách tin nhắn cho Admin
	r.GET("/api/admin/messages", func(c *gin.Context) {
	// Trong thực tế bạn nên kiểm tra quyền Admin ở đây bằng Middleware
		rows, err := db.Query("SELECT username, content, created_at FROM messages ORDER BY created_at DESC")
		if err != nil {
			c.JSON(500, gin.H{"error": "Không thể lấy tin nhắn"})
			return
		}
		defer rows.Close()

		type Msg struct {
			Username  string `json:"username"`
			Content   string `json:"content"`
			CreatedAt string `json:"created_at"`
		}
		var msgs []Msg
		for rows.Next() {
			var m Msg
			rows.Scan(&m.Username, &m.Content, &m.CreatedAt)
			msgs = append(msgs, m)
		}
		c.JSON(200, msgs)
	})

	r.GET("/admin", func(c *gin.Context) {
    c.File("./web/admin.html")
	})

	// 1. Khởi tạo Hub quản lý Chat
	hub := chat.NewHub()
	go hub.Run()

	// 2. Thiết lập endpoint WebSocket
	r.GET("/ws", func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing token"})
			return
		}

		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			log.Println("❌ Unauthorized WS:", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		log.Printf("✅ WS connected: %s (%s)",
			claims.Username,
			claims.UserID,
		)

		// Cho phép nâng cấp WebSocket
		chat.ServeWs(hub, c.Writer, c.Request)
	})


	// =====================
	// START SERVER
	// =====================
	fmt.Println("🚀 Server running at http://localhost:8080")
	r.Run(":8080")
}