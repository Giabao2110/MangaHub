package udp

import (
	"encoding/json"
	"log"
	"net"
	"sync"
)

// Cấu trúc tin nhắn thông báo
type Notification struct {
	Type      string `json:"type"`     // Ví dụ: "new_chapter"
	MangaID   string `json:"manga_id"` // Ví dụ: "OnePiece"
	Message   string `json:"message"`  // Ví dụ: "Chapter 1100 is out!"
	Timestamp int64  `json:"timestamp"`
}

type NotificationServer struct {
	Port    string
	Clients map[string]*net.UDPAddr // Danh sách địa chỉ các client đăng ký
	Lock    sync.Mutex
}

func NewServer(port string) *NotificationServer {
	return &NotificationServer{
		Port:    port,
		Clients: make(map[string]*net.UDPAddr),
	}
}

func (s *NotificationServer) Start() {
	// 1. Tạo địa chỉ UDP
	addr, err := net.ResolveUDPAddr("udp", ":"+s.Port)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Mở kết nối UDP (ListenUDP)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("📡 UDP Notification Server listening on port %s", s.Port)

	// 3. Vòng lặp lắng nghe tin nhắn từ Client (để đăng ký)
	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}

		// Khi client gửi bất kỳ tin gì đến, ta coi như họ muốn ĐĂNG KÝ nhận thông báo
		msg := string(buffer[:n])
		log.Printf("Received from %s: %s", clientAddr, msg)

		s.Lock.Lock()
		// Lưu địa chỉ client vào danh sách
		s.Clients[clientAddr.String()] = clientAddr
		s.Lock.Unlock()

		// Gửi lại tin xác nhận
		reply := []byte("✅ Subscribed to notifications!")
		conn.WriteToUDP(reply, clientAddr)
	}
}

// Hàm này dùng để Admin bắn thông báo cho toàn bộ Client
func (s *NotificationServer) Broadcast(notif Notification) {
	// Tạo kết nối tạm để gửi tin đi
	conn, _ := net.ListenPacket("udp", ":0")
	defer conn.Close()

	data, _ := json.Marshal(notif)

	s.Lock.Lock()
	defer s.Lock.Unlock()

	// Duyệt qua danh sách client và bắn tin đi (Fire and Forget)
	for _, addr := range s.Clients {
		conn.(*net.UDPConn).WriteToUDP(data, addr)
	}
	log.Printf("📢 Broadcasted to %d clients", len(s.Clients))
}
