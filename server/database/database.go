package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB - Biến global lưu connection đến database
// Được sử dụng ở nhiều nơi trong ứng dụng để query/update data
var DB *gorm.DB

// Connect - Hàm kết nối đến MySQL database
// Đọc cấu hình từ biến môi trường (.env file)
// Tạo DSN (Data Source Name) và kết nối qua GORM
func Connect() {
	var err error

	// Đọc thông tin kết nối từ biến môi trường
	dbUser := os.Getenv("DB_USER")         // Username MySQL (ví dụ: admin)
	dbPassword := os.Getenv("DB_PASSWORD") // Password MySQL
	dbHost := os.Getenv("DB_HOST")         // Host (ví dụ: localhost)
	dbPort := os.Getenv("DB_PORT")         // Port (ví dụ: 3306)
	dbName := os.Getenv("DB_NAME")         // Tên database (ví dụ: travel_db)

	// Tạo DSN (Data Source Name) theo định dạng MySQL
	// Format: user:password@tcp(host:port)/dbname?params
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Mở kết nối đến MySQL thông qua GORM
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// Nếu kết nối thất bại, dừng chương trình
		log.Fatal("Không thể kết nối database:", err)
	}

	log.Println("✅ Kết nối MySQL thành công!")
}
