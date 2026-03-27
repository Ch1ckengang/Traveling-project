package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB - Biến global lưu connection đến database
var DB *gorm.DB

// Connect - Kết nối đến SQLite database
// SQLite lưu dữ liệu trong file local — không cần cài server, không cần .env
// File travel.db sẽ tự động được tạo trong thư mục chạy server
func Connect() {
	var err error

	// Mở kết nối SQLite — chỉ cần tên file là xong
	DB, err = gorm.Open(sqlite.Open("travel.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Không thể kết nối database:", err)
	}

	log.Println("✅ Kết nối SQLite thành công! (file: travel.db)")
}
