package chat

import (
	"database/sql"
	"log"
)

// DB global cho chat
var chatDB *sql.DB

// InitChatDB được gọi từ main
func InitChatDB(db *sql.DB) {
	chatDB = db
}

// SaveMessage lưu chat vào SQLite
func SaveMessage(msg Message, userID string) {
	if chatDB == nil {
		log.Println("❌ Chat DB not initialized")
		return
	}

	_, err := chatDB.Exec(
		`INSERT INTO chat_messages (user_id, username, content) VALUES (?, ?, ?)`,
		userID,
		msg.Username,
		msg.Content,
	)

	if err != nil {
		log.Println("❌ Save chat error:", err)
	}

	log.Println("✅ Chat message saved for user:", userID)
}
	