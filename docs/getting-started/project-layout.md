# Project Layout & Architecture

Tài liệu này mô tả cấu trúc thư mục và kiến trúc của `chorus-controller` để bạn có thể điều hướng và phát triển nhanh chóng.

## Thư mục gốc

- `cmd/`: Điểm vào ứng dụng (entrypoint)
  - `main.go`: Khởi tạo cấu hình, DI các handler/service, tạo HTTP server và chạy
- `internal/`: Business logic và HTTP server (không export ra ngoài module)
  - `server/`: Khởi tạo router, middleware, đăng ký routes, lifecycle (Initialize/Run)
  - `handler/`: Lớp HTTP handler (Gin). Nhận `gin.Context`, validate input, gọi service
  - `service/`: Lớp nghiệp vụ. Chứa logic chính, validate domain, orchestrate repository
  - `repository/`: Truy cập dữ liệu (GORM). CRUD với Postgres, mapping model <-> DB
  - `domain/`: Khai báo domain models, DTOs, interfaces giữa các layer
  - `middleware/`: Các middleware dùng chung (error, auth token, system token)
  - `config/`: Đọc và biểu diễn cấu hình (env → struct)
  - `db/`: Khởi tạo kết nối GORM, tiện ích DB
  - `crypto/`: Mã hóa/giải mã, các helper liên quan đến secrets
  - `errors/`: Chuẩn hóa lỗi API (mã lỗi, HTTP code, thông điệp)
- `migrations/`: Các file migration SQL được quản lý bởi Atlas
- `docs/`: Tài liệu sử dụng và phát triển (getting-started, swagger, migrations, …)
- `scripts/`: Script tiện ích (setup env, apply migration, …)
- `build/`: Output binary khi `make build`
- `Makefile`: Tập hợp lệnh thường dùng (build/run/test, migrate, swagger, …)
- `atlas.hcl`: Cấu hình Atlas (migrations), môi trường local/dev
- `env.example`: Mẫu cấu hình môi trường
- `docker-compose.yml`: Compose cho Postgres và controller (tùy chọn)

## Dòng chảy request (HTTP)

```
Client → handler (HTTP) → service (business) → repository (DB) → Postgres
                              ↓
                         chorus-worker (gRPC)  [khi tạo job]
```

- Handler: chỉ làm nhiệm vụ HTTP (parse input, trả JSON). Không chứa business logic
- Service: xử lý nghiệp vụ, kiểm tra ràng buộc domain, gọi repository/worker
- Repository: tương tác cơ sở dữ liệu bằng GORM, trả domain models về service

## Authentication & Authorization

- Token chuẩn: dùng header `Authorization: Token <TOKEN>`
- System token:
  - Được tạo/đảm bảo tồn tại khi server khởi động (`Server.Initialize → EnsureSystemToken`)
  - Không bao giờ hết hạn (`expires_at = null`, JWT không có `exp`)
  - Có quyền truy cập các API quản trị: `GET /auth/tokens/detailed`, `POST /auth/revoke?token_id=...`, `DELETE /auth/tokens/:id`
- Middleware:
  - `middleware.AuthMiddleware`: xác thực token chuẩn cho các endpoint ghi dữ liệu
  - `middleware.SystemTokenAuth`: yêu cầu system token cho các endpoint quản trị
  - `middleware.ErrorHandler`: chuẩn hóa panic/recovery và format lỗi API

## Token Management (tóm tắt)

- Tạo token: `POST /auth/token` (có thể truyền `expires_at`; nếu không sẽ mặc định 24h; riêng token tên `system` sẽ không hết hạn)
- Liệt kê (public): `GET /auth/tokens` — trả về system token (kèm giá trị để sử dụng)
- Liệt kê chi tiết (system): `GET /auth/tokens/detailed` — trả tất cả tokens kèm JWT values
- Revoke theo ID (system): `POST /auth/revoke?token_id=<id>`
- Xóa theo ID (system): `DELETE /auth/tokens/:id` (không xóa được system token)

## Migrations & Schema

- Sử dụng Atlas để quản lý schema. Các models GORM trong `internal/domain/` là single source of truth
- Sinh migration từ models:
  - `make migrate-new NAME=<tên_thay_đổi>` → tool `cmd/schema/main.go` xuất `schema.sql` → Atlas diff → tạo file SQL trong `migrations/`
- Áp dụng migrations:
  - `make migrate-up`, rollback: `make migrate-down`, trạng thái: `make migrate-status`, checksum: `make migrate-hash`

## Liên kết nhanh

- Getting Started:
  - [Environment Setup](environment-setup.md)
  - [Development Guide](development.md)
  - [Migrations](migrations.md)
  - [Authentication](authentication.md)
  - [Swagger Authentication](swagger-authentication.md)

Nếu bạn mới tham gia dự án, hãy đọc lần lượt Environment Setup → Development → Migrations → Authentication để nắm luồng phát triển chuẩn.
