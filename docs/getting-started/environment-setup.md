# Environment Setup Guide

Hướng dẫn thiết lập môi trường cho Chorus Controller với file `.env` để bảo mật các thông tin cấu hình.

## Tổng quan

Chorus Controller sử dụng file `.env` để quản lý các thông tin cấu hình nhạy cảm như:
- Database connection strings
- JWT secrets
- API keys
- Environment-specific settings

## Quick Start

### 1. Setup Environment (Tự động)

```bash
# Chạy script setup tự động
make setup-env
```

Script sẽ:
- Tạo file `.env` từ template
- Generate JWT secret ngẫu nhiên
- Hỏi các thông tin cấu hình cần thiết
- Tạo file `.env` hoàn chỉnh

### 2. Setup Environment (Thủ công)

```bash
# Copy template
cp env.example .env

# Edit file .env với editor yêu thích
nano .env
# hoặc
vim .env
# hoặc
code .env
```

## Cấu hình Environment Variables

### Database Configuration

```bash
# PostgreSQL connection string
POSTGRES_DSN=postgres://username:password@host:port/database?sslmode=disable
```

**Ví dụ:**
- Local development: `postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable`
- Production: `postgres://user:secure_password@db.example.com:5432/chorus_prod?sslmode=require`

### HTTP Server Configuration

```bash
# Port cho HTTP server
HTTP_PORT=8081
```

### Worker gRPC Configuration

```bash
# Address của Chorus Worker gRPC service
WORKER_GRPC_ADDR=localhost:9670
```

### JWT Authentication Configuration

```bash
# Secret key cho JWT tokens (QUAN TRỌNG: Phải thay đổi trong production)
JWT_SECRET=your-secret-key-change-this-in-production-make-it-long-and-random

# Thời gian hết hạn của JWT tokens
JWT_EXPIRY=24h
```

**Lưu ý về JWT_SECRET:**
- Phải là chuỗi dài và ngẫu nhiên
- Không được chia sẻ hoặc commit vào git
- Thay đổi định kỳ trong production
- Sử dụng ít nhất 32 ký tự

### Environment Type

```bash
# Môi trường chạy (development/production)
ENV=development
```

## File .env Example

```bash
# Chorus Controller Configuration
# Database Configuration
POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable

# HTTP Server Configuration
HTTP_PORT=8081

# Worker gRPC Configuration
WORKER_GRPC_ADDR=localhost:9670

# JWT Authentication Configuration
JWT_SECRET=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6
JWT_EXPIRY=24h

# Development/Production Environment
ENV=development
```

## Security Best Practices

### 1. File Permissions

```bash
# Đặt quyền chỉ đọc cho owner
chmod 600 .env
```

### 2. Git Ignore

File `.env` đã được thêm vào `.gitignore`:
```gitignore
# Environment files
.env
.env.local
.env.*.local
```

### 3. Production Security

- ✅ **Sử dụng HTTPS** trong production
- ✅ **Thay đổi JWT_SECRET** từ default
- ✅ **Sử dụng strong passwords** cho database
- ✅ **Restrict database access** bằng firewall
- ✅ **Monitor access logs** thường xuyên
- ✅ **Rotate secrets** định kỳ

### 4. Environment Separation

Tạo các file riêng cho từng môi trường:

```bash
# Development
.env.development

# Staging  
.env.staging

# Production
.env.production
```

## Commands

### Setup Commands

```bash
# Setup environment tự động
make setup-env

# Kiểm tra environment
make env-check

# Chạy application (tự động check env)
make run
```

### Manual Commands

```bash
# Load environment variables
source .env

# Check environment variables
env | grep -E '^(POSTGRES_DSN|HTTP_PORT|WORKER_GRPC_ADDR|JWT_SECRET|JWT_EXPIRY|ENV)='

# Generate random JWT secret
openssl rand -base64 32
```

## Troubleshooting

### Lỗi "Could not load .env file"

```bash
# Kiểm tra file .env có tồn tại
ls -la .env

# Kiểm tra quyền file
ls -la .env

# Tạo lại file .env
make setup-env
```

### Lỗi Database Connection

```bash
# Kiểm tra PostgreSQL đang chạy
pg_isready -h localhost -p 5432

# Test connection
psql "postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable"

# Kiểm tra POSTGRES_DSN trong .env
grep POSTGRES_DSN .env
```

### Lỗi JWT Secret

```bash
# Kiểm tra JWT_SECRET có được set
grep JWT_SECRET .env

# Generate secret mới
openssl rand -base64 32

# Update trong .env
sed -i 's/JWT_SECRET=.*/JWT_SECRET=new_secret_here/' .env
```

### Lỗi Port Already in Use

```bash
# Kiểm tra port đang được sử dụng
lsof -i :8081

# Kill process sử dụng port
kill -9 $(lsof -t -i:8081)

# Hoặc thay đổi port trong .env
sed -i 's/HTTP_PORT=8081/HTTP_PORT=8082/' .env
```

## Development vs Production

### Development

```bash
# .env cho development
POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable
HTTP_PORT=8081
WORKER_GRPC_ADDR=localhost:9670
JWT_SECRET=dev-secret-key-not-for-production
JWT_EXPIRY=24h
ENV=development
```

### Production

```bash
# .env cho production
POSTGRES_DSN=postgres://prod_user:secure_password@prod-db.example.com:5432/chorus_prod?sslmode=require
HTTP_PORT=80
WORKER_GRPC_ADDR=worker.example.com:9670
JWT_SECRET=very-long-and-random-production-secret-key-32-chars-minimum
JWT_EXPIRY=1h
ENV=production
```

## Docker Integration

### Docker Compose với .env

```yaml
# docker-compose.yml
version: '3.8'
services:
  chorus-controller:
    build: .
    ports:
      - "${HTTP_PORT:-8081}:${HTTP_PORT:-8081}"
    environment:
      - POSTGRES_DSN=${POSTGRES_DSN}
      - JWT_SECRET=${JWT_SECRET}
      - WORKER_GRPC_ADDR=${WORKER_GRPC_ADDR}
    env_file:
      - .env
```

### Docker Run với .env

```bash
# Chạy với .env file
docker run --env-file .env chorus-controller

# Hoặc với environment variables
docker run \
  -e POSTGRES_DSN="postgres://..." \
  -e JWT_SECRET="..." \
  -e HTTP_PORT="8081" \
  chorus-controller
```

## Monitoring và Logging

### Environment Variables trong Logs

```bash
# Log environment info (không log sensitive data)
echo "Environment: $ENV"
echo "HTTP Port: $HTTP_PORT"
echo "Worker gRPC: $WORKER_GRPC_ADDR"
# Không log JWT_SECRET hoặc database password
```

### Health Check với Environment

```bash
# Health check endpoint
curl http://localhost:8081/health

# Check environment
curl -H "Authorization: Bearer $TOKEN" http://localhost:8081/auth/tokens
```
