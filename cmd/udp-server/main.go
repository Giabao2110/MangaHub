package main

import (
	"log"
	"mangahub/internal/udp"
	"net"
	"net/rpc" // Đây là RPC chuẩn của hệ thống
	"time"

	// Đặt biệt danh là "shared" để tránh trùng tên
	shared "mangahub/internal/rpc"
)

// Tạo một struct để hứng các lệnh RPC từ Admin
type AdminService struct {
	UDPServer *udp.NotificationServer
}

// [SỬA LẠI] Dùng "shared" thay vì "rpc" cho các tham số
func (s *AdminService) TriggerBroadcast(args *shared.BroadcastArgs, reply *shared.BroadcastReply) error {
	log.Printf("👮 Admin ordered broadcast: %s - %s", args.MangaID, args.Message)

	// Tạo thông báo
	notif := udp.Notification{
		Type:      "admin_alert",
		MangaID:   args.MangaID,
		Message:   args.Message,
		Timestamp: time.Now().Unix(),
	}

	// Ra lệnh cho UDP Server bắn tin đi
	s.UDPServer.Broadcast(notif)

	// Trả kết quả về cho Admin
	reply.Status = "Success"
	reply.Count = len(s.UDPServer.Clients)
	return nil
}

func main() {
	// 1. Khởi tạo UDP Server (Port 9091)
	udpServer := udp.NewServer("9091")

	// 2. Setup RPC Server (Cổng sau dành cho Admin - Port 1234)
	adminService := &AdminService{UDPServer: udpServer}

	// Đăng ký dịch vụ với thư viện chuẩn rpc
	rpc.Register(adminService)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("RPC Listen error:", err)
	}

	log.Println("✅ Admin RPC Interface listening on port 1234")

	// Chạy RPC ở một luồng riêng (Goroutine)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	// 3. Chạy UDP Server chính
	udpServer.Start()
}
