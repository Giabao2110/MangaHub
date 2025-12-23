package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = []byte("mangahub_secret_2025")

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./mangahub.db")
	if err != nil {
		log.Fatal("Lỗi kết nối DB:", err)
	}

	// Tạo bảng (Sử dụng '=' thay vì ':=' để tránh lỗi trùng biến)
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		role TEXT DEFAULT 'user',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS favorites (user_id INTEGER, manga_slug TEXT, PRIMARY KEY(user_id, manga_slug));
	CREATE TABLE IF NOT EXISTS comments (id INTEGER PRIMARY KEY, user_id INTEGER, manga_slug TEXT, content TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("Lỗi tạo bảng: %q\n", err)
	}

	// Tạo Admin mặc định
	hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), 14)
	// Dùng IGNORE để không lỗi nếu admin đã tồn tại
	db.Exec("INSERT OR IGNORE INTO users (username, password, role) VALUES (?, ?, ?)", "admin", string(hash), "admin")

	return db
}

func main() {
	db := InitDB()

	// Khởi tạo Router
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// --- CẤU HÌNH FILE TĨNH ---
	r.Static("/static", "./web") // Chứa style.css, app.js
	r.Static("/data", "./data")  // Chứa ảnh truyện

	// Routing cho các trang HTML
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/login", "./web/login.html")
	r.StaticFile("/register", "./web/register.html")

	// --- API AUTH ---

	// 1. Đăng Ký
	r.POST("/api/register", func(c *gin.Context) {
		var u struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
		_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", u.Username, string(hash))

		if err != nil {
			c.JSON(400, gin.H{"error": "Tên đăng nhập đã tồn tại!"})
			return
		}
		c.JSON(200, gin.H{"message": "Đăng ký thành công!"})
	})

	// 2. Đăng Nhập
	r.POST("/api/login", func(c *gin.Context) {
		var u struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
			return
		}

		var id int
		var hash, role string

		// Khai báo err rõ ràng để tránh lỗi 'no new variables'
		err := db.QueryRow("SELECT id, password, role FROM users WHERE username=?", u.Username).Scan(&id, &hash, &role)

		if err != nil {
			c.JSON(401, gin.H{"error": "Tài khoản không tồn tại"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password)); err != nil {
			c.JSON(401, gin.H{"error": "Sai mật khẩu"})
			return
		}

		// Tạo Token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id": id, "user": u.Username, "role": role, "exp": time.Now().Add(24 * time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString(SECRET_KEY)

		c.JSON(200, gin.H{"token": tokenString, "role": role, "username": u.Username})
	})

	// --- API MANGAS ---
	r.GET("/api/mangas", func(c *gin.Context) {
		var mangas []gin.H
		files, _ := ioutil.ReadDir("./data")
		for _, f := range files {
			if f.IsDir() {
				mangas = append(mangas, gin.H{
					"slug":     f.Name(),
					"title":    strings.Title(strings.ReplaceAll(f.Name(), "-", " ")),
					"cover":    fmt.Sprintf("/data/%s/cover.jpg", f.Name()),
					"category": "Manga",
				})
			}
		}
		c.JSON(200, mangas)
	})

	r.GET("/api/read/:slug/:chap", func(c *gin.Context) {
		path := fmt.Sprintf("./data/%s/%s", c.Param("slug"), c.Param("chap"))
		files, err := ioutil.ReadDir(path)
		if err != nil {
			c.JSON(404, gin.H{"error": "Không tìm thấy chapter"})
			return
		}
		var images []string
		for _, f := range files {
			if !f.IsDir() && (strings.HasSuffix(f.Name(), ".jpg") || strings.HasSuffix(f.Name(), ".png")) {
				images = append(images, fmt.Sprintf("/%s/%s", path, f.Name()))
			}
		}
		sort.Strings(images)
		c.JSON(200, images)
	})

	fmt.Println("🚀 Server đang chạy tại: http://localhost:8080")
	r.Run(":8080")
}
