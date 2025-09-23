# Swagger UI Authentication Guide

Hướng dẫn sử dụng Swagger UI với authentication cho Chorus Controller API.

## 🔐 Authentication trong Swagger UI

### 1. Truy cập Swagger UI

Mở trình duyệt và truy cập: `http://localhost:8081/swagger/index.html`

### 2. Authentication Policy

Chorus Controller sử dụng chính sách authentication phân biệt:

#### ✅ **Public Endpoints (Không cần token)**
- `GET /health` - Health check
- `POST /auth/token` - Tạo token mới
- `GET /storages` - Danh sách storages
- `GET /buckets` - Danh sách buckets
- `GET /storages/db` - Danh sách storages từ DB
- `GET /storages/{id}` - Chi tiết storage
- `GET /replications` - Danh sách replication jobs

#### 🔒 **Protected Endpoints (Cần token)**
- `GET /auth/tokens` - Danh sách tokens
- `POST /auth/revoke` - Revoke token
- `DELETE /auth/tokens/{id}` - Xóa token
- `POST /storages` - Tạo storage mới
- `PUT /storages/{id}` - Cập nhật storage
- `DELETE /storages/{id}` - Xóa storage
- `POST /replications` - Tạo replication job
- `POST /replications/pause` - Tạm dừng replication
- `POST /replications/resume` - Tiếp tục replication
- `DELETE /replications` - Xóa replication
- `POST /replications/switch/zero-downtime` - Switch zero downtime

## 🚀 Cách sử dụng Authentication trong Swagger UI

### Bước 1: Tạo Token

1. Mở Swagger UI
2. Tìm endpoint `POST /auth/token`
3. Click "Try it out"
4. Nhập thông tin:
   ```json
   {
     "name": "swagger-test-token",
     "description": "Token for Swagger UI testing"
   }
   ```
5. Click "Execute"
6. Copy token từ response (trường `token`)

### Bước 2: Cấu hình Authentication

1. Click nút **"Authorize"** ở góc trên bên phải Swagger UI
2. Trong hộp "Value", nhập: `Token <your-token>`
   - Ví dụ: `Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
3. Click **"Authorize"**
4. Click **"Close"**

### Bước 3: Test Protected Endpoints

Sau khi authorize, bạn có thể:
- Test các protected endpoints
- Swagger UI sẽ tự động thêm `Authorization: Token <token>` vào headers
- Xem response với dữ liệu đã được decrypt

## 📋 Ví dụ sử dụng

### Tạo Storage mới (Protected)

1. Tìm `POST /storages`
2. Click "Try it out"
3. Nhập dữ liệu:
   ```json
   {
     "name": "my-s3-storage",
     "address": "s3.amazonaws.com",
     "provider": "aws",
     "user": "myuser",
     "access_key": "AKIAIOSFODNN7EXAMPLE",
     "secret_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
   }
   ```
4. Click "Execute"
5. Kiểm tra response: `201 Created`

### Xem danh sách Storages (Public)

1. Tìm `GET /storages/db`
2. Click "Try it out"
3. Click "Execute"
4. Xem danh sách storages với secret keys đã được decrypt

### Tạo Replication Job (Protected)

1. Tìm `POST /replications`
2. Click "Try it out"
3. Nhập dữ liệu:
   ```json
   {
     "user": "myuser",
     "from": "source-storage",
     "to": "destination-storage",
     "buckets": ["bucket1", "bucket2"]
   }
   ```
4. Click "Execute"

## 🔧 Troubleshooting

### Lỗi "Authorization header is required"

**Nguyên nhân**: Chưa cấu hình token hoặc token không hợp lệ

**Giải pháp**:
1. Kiểm tra đã click "Authorize" chưa
2. Kiểm tra format token: `Bearer <token>` (có dấu cách)
3. Tạo token mới nếu token cũ hết hạn

### Lỗi "Invalid or expired token"

**Nguyên nhân**: Token đã hết hạn hoặc không hợp lệ

**Giải pháp**:
1. Tạo token mới bằng `POST /auth/token`
2. Cập nhật token trong Swagger UI
3. Kiểm tra JWT_EXPIRY trong .env

### Lỗi "Token not found"

**Nguyên nhân**: Token không tồn tại trong database

**Giải pháp**:
1. Kiểm tra database connection
2. Tạo token mới
3. Kiểm tra logs application

## 🎯 Best Practices

### 1. Token Management
- Tạo token với tên mô tả rõ ràng
- Sử dụng token khác nhau cho môi trường dev/prod
- Revoke token khi không sử dụng

### 2. Security
- Không chia sẻ token trong code hoặc logs
- Sử dụng HTTPS trong production
- Thay đổi JWT_SECRET định kỳ

### 3. Testing
- Test cả public và protected endpoints
- Kiểm tra encryption/decryption hoạt động
- Verify error handling

## 📚 API Reference

### Authentication Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| POST | `/auth/token` | ❌ | Tạo token mới |
| GET | `/auth/tokens` | ❌ | Trả về system token |
| GET | `/auth/tokens/detailed` | ✅ (system) | Danh sách tất cả tokens kèm giá trị |
| POST | `/auth/revoke?token_id=<id>` | ✅ (system) | Revoke token theo ID |
| DELETE | `/auth/tokens/{id}` | ✅ (system) | Xóa token (trừ system) |

### Storage Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/storages` | ❌ | Danh sách storages (worker) |
| GET | `/storages/db` | ❌ | Danh sách storages (DB) |
| GET | `/storages/{id}` | ❌ | Chi tiết storage |
| POST | `/storages` | ✅ | Tạo storage mới |
| PUT | `/storages/{id}` | ✅ | Cập nhật storage |
| DELETE | `/storages/{id}` | ✅ | Xóa storage |

### Replication Endpoints

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/replications` | ❌ | Danh sách replication jobs |
| POST | `/replications` | ✅ | Tạo replication job |
| POST | `/replications/pause` | ✅ | Tạm dừng replication |
| POST | `/replications/resume` | ✅ | Tiếp tục replication |
| DELETE | `/replications` | ✅ | Xóa replication |
| POST | `/replications/switch/zero-downtime` | ✅ | Switch zero downtime |

## 🔗 Links

- [Swagger UI](http://localhost:8081/swagger/index.html)
- [Health Check](http://localhost:8081/health)
- [Environment Setup Guide](environment-setup.md)
- [Authentication Guide](authentication.md)
