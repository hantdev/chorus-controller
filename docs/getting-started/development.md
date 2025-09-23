### Tổng quan
`chorus-controller` là control plane quản lý các công việc nhân bản (replication) S3 do `chorus-worker` thực thi. Tài liệu này là bản rút gọn, tránh trùng lặp bằng cách trỏ tới các hướng dẫn chi tiết trong cùng mục Getting Started.

### Yêu cầu
- Go 1.25+
- Docker/Docker Compose (Postgres dev)
- Atlas CLI (migrations)
- Dịch vụ gRPC `chorus-worker` (mặc định `localhost:9670`)

### Thiết lập môi trường (tóm tắt)
- Tạo và cấu hình `.env`: xem chi tiết ở `environment-setup.md`
- Mẫu `.env`: file `env.example` ở thư mục gốc

### Postgres cho môi trường phát triển
- Cách khởi chạy bằng Docker Compose/container, thiết lập DSN, kiểm tra kết nối: xem `environment-setup.md` (mục Database + Troubleshooting)

### Migrations (Atlas)
- Tổng quan, quy trình sinh/áp dụng migrations từ GORM models và các lệnh Makefile: xem `migrations.md`
- Lệnh nhanh:
```bash
make migrate-up
make migrate-down
make migrate-status
make migrate-new NAME=my_change
make migrate-hash
```

### Chạy ứng dụng (quick)
```bash
make run           # chạy server
# hoặc
go run ./cmd

# Docker Compose (controller + postgres)
docker compose up --build -d

# Docker image
make docker-build
make migrate-up
make docker-run
```

### Build & Test (quick)
```bash
make build
make test
# hoặc
go build -o bin/chorus-controller ./cmd
go test ./...
```

### API & Swagger
- Cách authorize, tiền tố header `Authorization: Token <TOKEN>` và thử endpoints: xem `swagger-authentication.md` và `authentication.md`
- Swagger UI: `http://localhost:8081/swagger/index.html`

### Tích hợp với chorus-worker
- Thiết lập `WORKER_GRPC_ADDR` và lưu ý khi chạy worker cục bộ: xem `environment-setup.md`

### Quy trình phát triển nhanh
```bash
# 1) Khởi chạy Postgres (xem environment-setup.md)
# 2) Áp dụng migrations: make migrate-up
# 3) Khởi động worker (nếu cần)
# 4) Chạy controller: make run
# 5) Mở Swagger: http://localhost:8081/swagger/index.html
```

### Xử lý sự cố (tham chiếu)
- Lỗi môi trường, Postgres, port, JWT secret: `environment-setup.md`
- Lỗi migrations/Atlas: `migrations.md`

### Tham chiếu nhanh Makefile
```bash
make build            # build binary
make run              # chạy server
make test             # chạy tests
make clean            # xóa bin/
make migrate-up       # áp dụng migrations
make migrate-down     # rollback gần nhất
make migrate-status   # trạng thái migrations
make migrate-new      # tạo migration mới
make migrate-hash     # cập nhật checksums
make docker-build     # build docker image
make docker-run       # chạy container (env mapping)
```

### Tài liệu liên quan
- Kiến trúc & cấu trúc thư mục: `project-layout.md`
- Thiết lập môi trường & lệnh: `environment-setup.md`
- Authentication & Token: `authentication.md`
- Swagger & Authorization header: `swagger-authentication.md`
- Migrations & Atlas: `migrations.md`
- Quy trình phát triển API mới: `api-development.md`
