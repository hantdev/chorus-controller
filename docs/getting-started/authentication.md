# Authentication System

Chorus Controller sử dụng hệ thống authentication đơn giản với JWT tokens để bảo vệ các API endpoints.

## Tổng quan

- **Không cần user/password**: Hệ thống chỉ cần tạo token để truy cập API
- **JWT Tokens**: Sử dụng JWT để xác thực và phân quyền
- **Database Storage**: Token information được lưu trong database
- **Middleware Protection**: Tất cả API endpoints được bảo vệ bằng authentication middleware

## Cách sử dụng

### 1. Tạo Token

Để tạo một token mới, gửi POST request đến `/auth/token`:

```bash
curl -X POST http://localhost:8081/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-api-client",
    "description": "Token for API access"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "name": "my-api-client",
  "expires_at": "2024-01-02T12:00:00Z",
  "created_at": "2024-01-01T12:00:00Z"
}
```

### 2. Sử dụng Token

Sử dụng token trong header `Authorization` với tiền tố `Token`:

```bash
curl -X GET http://localhost:8081/storages \
  -H "Authorization: Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 3. Quản lý Tokens

#### Liệt kê system token (public):
```bash
curl -X GET http://localhost:8081/auth/tokens
```

#### Liệt kê tất cả tokens kèm giá trị (yêu cầu system token):
```bash
curl -X GET http://localhost:8081/auth/tokens/detailed \
  -H "Authorization: Token <SYSTEM_TOKEN>"
```

#### Vô hiệu hóa token theo ID (yêu cầu system token):
```bash
curl -X POST "http://localhost:8081/auth/revoke?token_id=<TOKEN_ID>" \
  -H "Authorization: Token <SYSTEM_TOKEN>"
```

#### Xóa token:
```bash
curl -X DELETE http://localhost:8081/auth/tokens/<TOKEN_ID> \
  -H "Authorization: Token <SYSTEM_TOKEN>"
```

## API Endpoints

### Public Endpoints (Không cần authentication)
- `GET /health` — Health check
- `POST /auth/token` — Tạo token mới
- `GET /auth/tokens` — Trả về system token (để sử dụng)
- `GET /storages`, `GET /storages/db`, `GET /storages/:id`
- `GET /buckets`
- `GET /replications`

### Protected Endpoints
- System token required (`Authorization: Token <SYSTEM_TOKEN>`):
  - `GET /auth/tokens/detailed` — Danh sách tất cả tokens kèm giá trị
  - `POST /auth/revoke?token_id=<id>` — Vô hiệu hóa token theo ID
  - `DELETE /auth/tokens/:id` — Xóa token (không xóa được system token)
- Normal token required (`Authorization: Token <TOKEN>`):
  - `POST /storages`, `PUT /storages/:id`, `DELETE /storages/:id`
  - `POST /replications`, `POST /replications/pause`, `POST /replications/resume`, `DELETE /replications`, `POST /replications/switch/zero-downtime`

## Cấu hình

### Environment Variables

```bash
# JWT Secret (thay đổi trong production)
JWT_SECRET=your-secret-key-change-this-in-production

# JWT Expiry (default: 24h)
JWT_EXPIRY=24h

# Database connection
POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/chorus?sslmode=disable

# HTTP Port
HTTP_PORT=8081

# Worker gRPC address
WORKER_GRPC_ADDR=localhost:9670
```

### Database Migration

Chạy migrations bằng Atlas:

```bash
make migrate-up          # apply tất cả migrations
make migrate-status      # kiểm tra trạng thái
make migrate-down        # rollback migration gần nhất
```

## Bảo mật

### Production Checklist

- ✅ **Thay đổi JWT_SECRET**: Sử dụng secret key mạnh và unique
- ✅ **HTTPS**: Sử dụng HTTPS trong production
- ✅ **Token Expiry**: Đặt thời gian hết hạn hợp lý
- ✅ **Database Security**: Bảo mật database connection
- ✅ **Monitoring**: Monitor token usage và suspicious activities

### Token Security

- Tokens được hash trước khi lưu vào database
- JWT tokens chứa thông tin tối thiểu cần thiết
- Tokens có thể bị vô hiệu hóa hoặc xóa bất cứ lúc nào
- Mỗi token có thời gian hết hạn
