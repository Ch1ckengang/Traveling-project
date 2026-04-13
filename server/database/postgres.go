package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB - Biến global lưu connection đến database
var DB *gorm.DB

// Connect - Kết nối đến PostgreSQL database
// Có thể cấu hình qua biến môi trường hoặc dùng default cho local.
func Connect() {
	var err error

	// Load .env when running locally; ignore if file does not exist.
	_ = godotenv.Load()

	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := getenv("DB_USER", "postgres")
	password := getenv("DB_PASSWORD", "123456")
	name := getenv("DB_NAME", "travel_db")
	sslMode := getenv("DB_SSLMODE", "disable")
	timezone := getenv("DB_TIMEZONE", "Asia/Ho_Chi_Minh")

	// Prefer full DSN when provided (useful for Docker/cloud deployments).
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			host,
			user,
			password,
			name,
			port,
			sslMode,
			timezone,
		)
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Không thể kết nối database:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Không thể lấy SQL DB instance:", err)
	}

	// Keep a sane default connection pool for API workloads.
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Không thể ping PostgreSQL:", err)
	}

	log.Printf("✅ Kết nối PostgreSQL thành công! (%s:%s/%s)", host, port, name)
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
