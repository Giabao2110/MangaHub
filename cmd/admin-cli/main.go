package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc" // RPC chuẩn
	"os"
	"strings"

	// Đặt biệt danh là "shared"
	shared "mangahub/internal/rpc"
)

func main() {
	// 1. Kết nối tới cổng quản trị của Server (Port 1234)
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("❌ Không thể kết nối tới Server Admin:", err)
	}
	fmt.Println("✅ Đã kết nối tới hệ thống quản trị MangaHub!")
	fmt.Println("Gõ tin nhắn để thông báo cho toàn bộ User (hoặc 'exit' để thoát)")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n(Admin) > Nhập nội dung thông báo: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			break
		}
		if text == "" {
			continue
		}

		// 2. Chuẩn bị dữ liệu gửi đi (Dùng shared)
		args := &shared.BroadcastArgs{
			Message: text,
			MangaID: "System",
		}
		var reply shared.BroadcastReply

		// 3. GỌI HÀM TỪ XA
		err = client.Call("AdminService.TriggerBroadcast", args, &reply)
		if err != nil {
			log.Println("❌ Lỗi RPC:", err)
			continue
		}

		fmt.Printf("🚀 Đã gửi thành công cho %d users!\n", reply.Count)
	}
}
