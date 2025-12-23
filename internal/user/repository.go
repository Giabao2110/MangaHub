package user

import (
	"database/sql"
	//"mangahub/pkg/utils" // Import utils để dùng hàm băm password nếu cần, hoặc xử lý ở handler
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// CreateUser tạo người dùng mới
func (r *Repository) CreateUser(user *User) error {
	// 1. Mã hóa mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 2. Lưu vào DB
	query := `INSERT INTO users (id, username, password_hash) VALUES (?, ?, ?)`
	_, err = r.DB.Exec(query, user.ID, user.Username, string(hashedPassword))
	return err
}

// GetUserByUsername tìm người dùng theo username để đăng nhập
func (r *Repository) GetUserByUsername(username string) (*User, error) {
	row := r.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username)
	
	user := &User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, err
	}
	return user, nil
}