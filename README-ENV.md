# Environment Configuration với .env

Chorus Controller đã được cập nhật để sử dụng file `.env` cho việc quản lý cấu hình một cách bảo mật.

## 🚀 Quick Start

### 1. Setup Environment

```bash
# Tự động setup (khuyến nghị)
make setup-env

# Hoặc thủ công
cp env.example .env
# Edit .env với editor yêu thích
```

### 2. Kiểm tra Configuration

```bash
# Kiểm tra .env file
make env-check

# Xem nội dung .env (ẩn sensitive data)
cat .env
```

### 3. Chạy Application

```bash
# Chạy với environment check tự động
make run

# Hoặc chạy trực tiếp
go run cmd/main.go
```

## 📋 Configuration Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `POSTGRES_DSN` | Database connection string | `postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable` | ✅ |
| `HTTP_PORT` | HTTP server port | `8081` | ✅ |
| `WORKER_GRPC_ADDR` | Chorus Worker gRPC address | `localhost:9670` | ✅ |
| `JWT_SECRET` | JWT signing secret | Random generated | ✅ |
| `JWT_EXPIRY` | JWT token expiry | `24h` | ✅ |
| `ENV` | Environment type | `development` | ❌ |

## 🔒 Security Features

### ✅ Implemented
- **File `.env` được ignore** trong `.gitignore`
- **JWT secret tự động generate** ngẫu nhiên
- **Environment variables** được load từ file `.env`
- **Fallback values** nếu không có `.env`
- **File permissions** được check trong script

### 🛡️ Best Practices
- Sử dụng JWT secret dài và ngẫu nhiên
- Không commit file `.env` vào git
- Sử dụng HTTPS trong production
- Thay đổi secrets định kỳ

## 📁 File Structure

```
chorus-controller/
├── .env                     # Environment variables (không commit)
├── env.example       # Template cho .env
├── scripts/
│   └── setup-env.sh        # Script setup tự động
├── docs/
│   ├── environment-setup.md # Hướng dẫn chi tiết
│   └── authentication.md   # Hướng dẫn authentication
└── internal/config/
    └── config.go           # Config loader với godotenv
```

## 🛠️ Commands

### Environment Management
```bash
# Setup environment tự động
make setup-env

# Kiểm tra environment
make env-check

# Chạy application
make run
```

### Manual Commands
```bash
# Generate JWT secret
openssl rand -base64 32

# Load environment
source .env

# Check variables
env | grep -E '^(POSTGRES_DSN|HTTP_PORT|WORKER_GRPC_ADDR|JWT_SECRET|JWT_EXPIRY|ENV)='

# Test database connection
psql "$POSTGRES_DSN"
```

## 🔧 Development Workflow

### 1. First Time Setup
```bash
# Clone repository
git clone <repository-url>
cd chorus-controller

# Setup environment
make setup-env

# Apply database migrations
make migrate-apply

# Start application
make run
```

### 2. Daily Development
```bash
# Check environment
make env-check

# Run application
make run

# Test API
curl http://localhost:8081/health
```

### 3. Production Deployment
```bash
# Create production .env
cp config.env.example .env.production

# Edit với production values
vim .env.production

# Set environment
export ENV=production

# Run application
make run
```

## 🐛 Troubleshooting

### Lỗi "Could not load .env file"
```bash
# Kiểm tra file tồn tại
ls -la .env

# Tạo lại file
make setup-env
```

### Lỗi Database Connection
```bash
# Kiểm tra PostgreSQL
pg_isready -h localhost -p 5432

# Test connection
psql "postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable"
```

### Lỗi Port Already in Use
```bash
# Tìm process sử dụng port
lsof -i :8081

# Kill process
kill -9 $(lsof -t -i:8081)

# Hoặc thay đổi port trong .env
sed -i 's/HTTP_PORT=8081/HTTP_PORT=8082/' .env
```

## 📚 Documentation

- [Environment Setup Guide](docs/environment-setup.md) - Hướng dẫn chi tiết
- [Authentication Guide](docs/authentication.md) - Hướng dẫn authentication
- [Migration Guide](docs/migrations.md) - Hướng dẫn database migrations

## 🎯 Next Steps

1. **Review configuration**: Kiểm tra file `.env` đã được tạo
2. **Apply migrations**: Chạy `make migrate-apply`
3. **Test authentication**: Tạo token và test API
4. **Deploy to production**: Sử dụng production values

## 💡 Tips

- Sử dụng `make setup-env` cho lần đầu setup
- Luôn kiểm tra `.env` trước khi deploy
- Backup `.env` file trong production
- Sử dụng environment-specific files (`.env.development`, `.env.production`)
- Monitor logs để detect configuration issues
