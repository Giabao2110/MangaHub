package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB khởi tạo kết nối và tự động tạo các bảng cần thiết
func InitDB(path string) *sql.DB {
	// Mở kết nối tới file sqlite3
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("❌ Không thể mở database: %v", err)
	}

	// Bật Foreign Keys để đảm bảo ràng buộc dữ liệu
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Println("⚠️ Lỗi bật Foreign Keys:", err)
	}

	// Danh sách các câu lệnh tạo bảng (Schemas)
	schemas := []string{
		// 1. Bảng Users (Lưu tài khoản và quyền hạn)
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE,
			password_hash TEXT,
			role TEXT DEFAULT 'user', 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 2. Bảng Wishlist (Lưu các bộ truyện yêu thích của User)
		`CREATE TABLE IF NOT EXISTS wishlist (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			manga_slug TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(username, manga_slug) -- Đảm bảo không bị trùng lặp yêu thích
		);`,

		// 3. Bảng Messages (Dùng làm Comment hoặc gửi tin nhắn cho Admin)
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 4. Bảng Mangas (Dùng để lưu metadata truyện nếu bạn muốn quản lý nâng cao hơn sau này)
		`CREATE TABLE IF NOT EXISTS mangas (
			id TEXT PRIMARY KEY,
			title TEXT,
			author TEXT,
			description TEXT,
			cover_url TEXT,
			category TEXT
		);`,
	}

	// Thực thi từng câu lệnh tạo bảng
	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			log.Printf("❌ Lỗi khi khởi tạo schema: %v", err)
		}
	}

	return db
}