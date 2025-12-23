package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Cấu trúc gói tin giống Server
type ProgressUpdate struct {
	UserID    string `json:"user_id"`
	MangaID   string `json:"manga_id"`
	Chapter   int    `json:"chapter"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	// 1. Kết nối đến TCP Server
	serverAddress := "localhost:9090"
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("❌ Không thể kết nối đến Server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("✅ Đã kết nối đến Server tại", serverAddress)
	fmt.Println("------------------------------------------------")

	// 2. [QUAN TRỌNG] Tạo luồng lắng nghe tin nhắn từ Server (Goroutine)
	// Để khi Terminal 3 gửi, Terminal 2 (mình) cũng nhận được ngay
	go func() {
		decoder := json.NewDecoder(conn)
		for {
			var update ProgressUpdate
			if err := decoder.Decode(&update); err != nil {
				fmt.Println("\n⚠️ Mất kết nối tới Server!")
				os.Exit(0)
			}
			// In ra thông báo khi nhận được broadcast
			fmt.Printf("\n🔔 [ĐỒNG BỘ] User '%s' đang đọc '%s' chap %d\n> ",
				update.UserID, update.MangaID, update.Chapter)
		}
	}()

	// 3. Luồng chính: Gửi dữ liệu lên Server
	reader := bufio.NewReader(os.Stdin)
	// Giả sử tên user là tên thư mục máy tính
	userID, _ := os.Hostname()

	for {
		fmt.Print("> Nhập tên truyện (ví dụ: OnePiece): ")
		mangaID, _ := reader.ReadString('\n')
		mangaID = strings.TrimSpace(mangaID)

		if mangaID == "exit" {
			break
		}

		fmt.Print("> Nhập số chapter (ví dụ: 100): ")
		chapStr, _ := reader.ReadString('\n')
		chapter, _ := strconv.Atoi(strings.TrimSpace(chapStr))

		// Đóng gói JSON
		update := ProgressUpdate{
			UserID:    userID,
			MangaID:   mangaID,
			Chapter:   chapter,
			Timestamp: time.Now().Unix(),
		}

		// Gửi đi
		encoder := json.NewEncoder(conn)
		if err := encoder.Encode(update); err != nil {
			fmt.Println("❌ Lỗi gửi dữ liệu:", err)
			break
		}
		// Chờ một chút để giao diện đẹp hơn
		time.Sleep(100 * time.Millisecond)
	}
}
