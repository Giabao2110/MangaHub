package rpc

// Payload gửi lên từ Admin
type BroadcastArgs struct {
	Message string // Nội dung thông báo muốn gửi
	MangaID string // Tên truyện
}

// Kết quả trả về cho Admin
type BroadcastReply struct {
	Status string // "OK" hoặc "Error"
	Count  int    // Số lượng client đã nhận được tin
}
