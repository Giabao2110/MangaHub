package tcp

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

// Cấu trúc gói tin
type ProgressUpdate struct {
	UserID    string `json:"user_id"`
	MangaID   string `json:"manga_id"`
	Chapter   int    `json:"chapter"`
	Timestamp int64  `json:"timestamp"`
}

// Struct quản lý Server
type ProgressSyncServer struct {
	Port        string
	Connections map[string]net.Conn
	Broadcast   chan ProgressUpdate
	Lock        sync.Mutex
}

// [QUAN TRỌNG] Hàm này chính là tcp.NewServer mà máy đang tìm kiếm
// Chữ N phải viết hoa
func NewServer(port string) *ProgressSyncServer {
	return &ProgressSyncServer{
		Port:        port,
		Connections: make(map[string]net.Conn),
		Broadcast:   make(chan ProgressUpdate),
	}
}

// Hàm khởi động Server
func (s *ProgressSyncServer) Start() {
	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
	defer listener.Close()

	// Chạy luồng broadcast
	go s.handleBroadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *ProgressSyncServer) handleConnection(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	log.Printf("New client connected: %s", addr)

	s.Lock.Lock()
	s.Connections[addr] = conn
	s.Lock.Unlock()

	defer func() {
		log.Printf("Client disconnected: %s", addr)
		s.Lock.Lock()
		delete(s.Connections, addr)
		s.Lock.Unlock()
		conn.Close()
	}()

	decoder := json.NewDecoder(conn)
	for {
		var update ProgressUpdate
		if err := decoder.Decode(&update); err != nil {
			break
		}

		log.Printf("Received: User %s read %s ch.%d", update.UserID, update.MangaID, update.Chapter)
		s.Broadcast <- update
	}
}

func (s *ProgressSyncServer) handleBroadcast() {
	for {
		msg := <-s.Broadcast
		data, _ := json.Marshal(msg)

		s.Lock.Lock()
		for addr, conn := range s.Connections {
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if _, err := conn.Write(data); err != nil {
				conn.Close()
				delete(s.Connections, addr)
			}
		}
		s.Lock.Unlock()
	}
}
