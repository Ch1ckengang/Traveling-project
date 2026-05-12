# 🚀 SETUP GUIDE - TRAVELING PROJECT

Hướng dẫn setup và chạy dự án Traveling từ đầu.

---

## 📋 YÊU CẦU HỆ THỐNG

### Backend
- **Go:** 1.25.0 hoặc mới hơn
- **PostgreSQL:** 14+ (hoặc Docker)
- **Git:** Để clone repository

### Frontend
- **Node.js:** 18+ hoặc 20+
- **npm:** 9+ (đi kèm với Node.js)

---

## 🗄️ SETUP DATABASE

### Option 1: Docker (Khuyến nghị)

```bash
# Start PostgreSQL với Docker Compose
cd server
docker-compose up -d

# Kiểm tra container đang chạy
docker ps

# Xem logs
docker-compose logs -f postgres
```

Database sẽ chạy trên:
- **Host:** localhost
- **Port:** 5432
- **Database:** travel_db
- **Username:** postgres
- **Password:** 123456

### Option 2: PostgreSQL Local

```bash
# Install PostgreSQL (Ubuntu/Debian)
sudo apt update
sudo apt install postgresql postgresql-contrib

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Tạo database
sudo -u postgres psql
CREATE DATABASE travel_db;
CREATE USER postgres WITH PASSWORD '123456';
GRANT ALL PRIVILEGES ON DATABASE travel_db TO postgres;
\q
```

---

## 🔧 SETUP BACKEND

### 1. Clone Repository

```bash
git clone <repository-url>
cd traveling-project
```

### 2. Setup Environment Variables

```bash
cd server
cp .env.example .env
```

Chỉnh sửa `.env` nếu cần:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123456
DB_NAME=travel_db

# JWT Secrets (CHANGE IN PRODUCTION!)
JWT_ACCESS_SECRET=your-secret-key-here
JWT_REFRESH_SECRET=your-refresh-secret-here

# Email Service (Optional - Dev mode works without this)
EMAIL_ENABLED=false
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

### 3. Install Dependencies

```bash
go mod download
go mod tidy
```

### 4. Run Database Migrations

Migrations chạy tự động khi start server lần đầu (GORM AutoMigrate).

### 5. Start Backend Server

```bash
# Development mode
go run cmd/server/main.go

# Build và run
go build -o bin/server cmd/server/main.go
./bin/server
```

Server sẽ chạy trên: **http://localhost:8080**

### 6. Verify Backend

```bash
# Test health check
curl http://localhost:8080/v1/api/tours

# Kết quả mong đợi: JSON với danh sách tours
```

---

## 🎨 SETUP FRONTEND

### 1. Navigate to Client Folder

```bash
cd client
```

### 2. Setup Environment Variables

```bash
cp .env.example .env
```

File `.env` mặc định:

```env
VITE_API_URL=http://localhost:8080/v1/api
```

### 3. Install Dependencies

```bash
npm install
```

### 4. Start Development Server

```bash
npm run dev
```

Frontend sẽ chạy trên: **http://localhost:5173**

### 5. Verify Frontend

Mở trình duyệt: http://localhost:5173

Bạn sẽ thấy trang chủ Traveling.

---

## 🧪 TESTING

### Test Backend API

```bash
# Test register
curl -X POST http://localhost:8080/v1/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }'

# Test send OTP
curl -X POST http://localhost:8080/v1/api/otp/send \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com"
  }'

# Check OTP in server logs
# OTP code sẽ được log ra console trong dev mode
```

### Test Frontend

1. **Register Flow:**
   - Vào http://localhost:5173/register
   - Điền form và submit
   - Check OTP trong server logs
   - Verify OTP

2. **Login Flow:**
   - Vào http://localhost:5173/login
   - Login với email đã verify
   - Check redirect về home page

3. **Browse Tours:**
   - Vào http://localhost:5173/tours
   - Xem danh sách tours
   - Click vào tour để xem chi tiết

---

## 📧 SETUP EMAIL SERVICE (OPTIONAL)

### Dev Mode (Default)

Email chỉ log ra console, không gửi thật.

```env
EMAIL_ENABLED=false
```

### Production Mode - Gmail

1. **Enable 2-Factor Authentication:**
   - Vào https://myaccount.google.com/security
   - Bật 2-Step Verification

2. **Generate App Password:**
   - Vào https://myaccount.google.com/apppasswords
   - Tạo app password mới
   - Copy 16-character password

3. **Update .env:**

```env
EMAIL_ENABLED=true
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-16-char-app-password
SMTP_FROM_EMAIL=your-email@gmail.com
SMTP_FROM_NAME=Traveling
```

### Production Mode - SendGrid

1. **Sign up:** https://sendgrid.com
2. **Create API Key:** Settings → API Keys → Create API Key
3. **Update .env:**

```env
EMAIL_ENABLED=true
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
SMTP_FROM_EMAIL=your-verified-sender@example.com
SMTP_FROM_NAME=Traveling
```

---

## 🐛 TROUBLESHOOTING

### Backend Issues

**Problem:** `cannot connect to database`
```bash
# Check PostgreSQL is running
docker ps  # for Docker
sudo systemctl status postgresql  # for local install

# Check connection
psql -h localhost -U postgres -d travel_db
```

**Problem:** `port 8080 already in use`
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

**Problem:** `go: module not found`
```bash
# Clean and reinstall
go clean -modcache
go mod download
go mod tidy
```

### Frontend Issues

**Problem:** `VITE_API_URL is not defined`
```bash
# Make sure .env exists
cp .env.example .env

# Restart dev server
npm run dev
```

**Problem:** `Cannot connect to backend`
```bash
# Check backend is running
curl http://localhost:8080/v1/api/tours

# Check CORS settings in server/cmd/server/main.go
# Should include: http://localhost:5173
```

**Problem:** `npm install fails`
```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

### Database Issues

**Problem:** `relation "users" does not exist`
```bash
# Restart server to trigger AutoMigrate
# Or manually run migrations
```

**Problem:** `password authentication failed`
```bash
# Check .env credentials match database
# Reset PostgreSQL password if needed
```

---

## 🔒 SECURITY NOTES

### Development

- Default credentials are for **development only**
- JWT secrets are weak - change in production
- CORS allows localhost - restrict in production

### Production Checklist

- [ ] Change all default passwords
- [ ] Generate strong JWT secrets (32+ characters)
- [ ] Enable HTTPS/SSL
- [ ] Configure proper CORS origins
- [ ] Enable email service
- [ ] Setup database backups
- [ ] Configure rate limiting
- [ ] Add monitoring and logging

---

## 📚 USEFUL COMMANDS

### Backend

```bash
# Run tests
go test ./...

# Format code
go fmt ./...

# Check for issues
go vet ./...

# Build for production
go build -o bin/server cmd/server/main.go

# Run with hot reload (install air first)
air
```

### Frontend

```bash
# Development
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint
```

### Database

```bash
# Connect to database
psql -h localhost -U postgres -d travel_db

# Backup database
pg_dump -h localhost -U postgres travel_db > backup.sql

# Restore database
psql -h localhost -U postgres travel_db < backup.sql

# Reset database
docker-compose down -v
docker-compose up -d
```

---

## 🎯 NEXT STEPS

Sau khi setup xong:

1. ✅ Test register và login flow
2. ✅ Browse tours
3. ✅ Create booking
4. 🔜 Setup payment gateway (NGÀY 2)
5. 🔜 Implement review module (NGÀY 3)
6. 🔜 Implement admin dashboard (NGÀY 5)

---

## 📞 SUPPORT

Nếu gặp vấn đề:

1. Check logs: server console và browser console
2. Verify environment variables
3. Check database connection
4. Review CHANGELOG-DAY1.md cho updates mới nhất

---

**Happy Coding! 🚀**
