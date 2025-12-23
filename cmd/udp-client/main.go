package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. Kết nối đến Server 9091
	serverAddr, _ := net.ResolveUDPAddr("udp", "localhost:9091")
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// 2. Gửi tin chào hỏi để Đăng Ký
	conn.Write([]byte("Hello Server!"))
	fmt.Println("📨 Đã gửi yêu cầu đăng ký...")

	// 3. Vòng lặp chờ tin nhắn từ Server
	buffer := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("\n🔔 THÔNG BÁO MỚI: %s\n", message)
	}
}
