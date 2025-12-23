package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	// Bật Foreign Keys
	db.Exec("PRAGMA foreign_keys = ON")

	schemas := []string{
		// 1. Bảng Users (Thêm cột Role để phân biệt Admin/User)
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE,
			password_hash TEXT,
			role TEXT DEFAULT 'user', 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		// 2. Bảng Mangas (Chứa thông tin truyện)
		`CREATE TABLE IF NOT EXISTS mangas (
			id TEXT PRIMARY KEY,
			title TEXT,
			author TEXT,
			description TEXT,
			cover_url TEXT,
			status TEXT, -- 'ongoing', 'completed'
			category TEXT
		);`,
		// 3. Bảng Chapters
		`CREATE TABLE IF NOT EXISTS chapters (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			manga_id TEXT,
			chapter_number INTEGER,
			title TEXT,
			content_images TEXT, -- JSON string chứa danh sách ảnh
			FOREIGN KEY(manga_id) REFERENCES mangas(id)
		);`,
		// 4. Bảng Favorites (Tủ truyện)
		`CREATE TABLE IF NOT EXISTS favorites (
			user_id TEXT,
			manga_id TEXT,
			PRIMARY KEY (user_id, manga_id)
		);`,
		// 5. Bảng Comments (Bình luận)
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT,
			manga_id TEXT,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			log.Printf("Error migrating schema: %v", err)
		}
	}

	return db
}
