package chat

import "log"

// Message cấu trúc tin nhắn trao đổi trong phòng chat
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Time     string `json:"time"`
}

// Hub quản lý các client và điều phối tin nhắn
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

// NewHub khởi tạo Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run chạy vòng lặp xử lý chat
func (h *Hub) Run() {
	log.Println("🚀 Chat Hub is running...")

	for {
		select {

		case client := <-h.register:
			h.clients[client] = true
			log.Printf("👤 User %s joined chat", client.Username)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("👋 User %s left chat", client.Username)
			}

		case message := <-h.broadcast:

			// 💾 LƯU CHAT VÀO SQLITE (CHỈ 1 LẦN)
			for client := range h.clients {
				SaveMessage(message, client.UserID)
				break
			}

			// 🔊 BROADCAST CHO TẤT CẢ CLIENT
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
