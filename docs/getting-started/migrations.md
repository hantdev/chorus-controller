# Database Migrations với Atlas

Hướng dẫn sử dụng Atlas để tự động sinh và quản lý database migrations từ GORM models.

## Cài đặt Atlas

```bash
# macOS
brew install atlas

# Linux/Windows
curl -sSf https://atlasgo.sh | sh
```

## Cấu hình

Hệ thống đã được cấu hình để tự động sinh migrations từ GORM models:
- `cmd/schema/main.go`: Tool để export schema từ GORM models
- `atlas.hcl`: Cấu hình Atlas environments
- `Makefile`: Các targets để quản lý migrations

## Sử dụng

### 1. Sinh migration từ model changes (Makefile)

Khi bạn thay đổi GORM models trong `internal/domain/models.go`, chạy:

```bash
# Tạo migration tự động từ GORM models
make migrate-new NAME=add_user_table

# Sau khi tạo migration, có thể cần chỉnh sửa file SQL trong migrations/
```

**Quy trình tự động:**
1. Tool sẽ đọc GORM models và sinh PostgreSQL schema
2. Atlas sẽ so sánh schema hiện tại với schema mong muốn
3. Tự động tạo file migration SQL với các thay đổi cần thiết

### 2. Apply migrations

```bash
make migrate-up

make migrate-down       # rollback 1 migration

# Hoặc sử dụng Atlas trực tiếp
GOWORK=off atlas migrate apply --env local
```

### 3. Kiểm tra trạng thái migrations

```bash
# Sử dụng Makefile
make migrate-status

# Hoặc sử dụng Atlas trực tiếp
GOWORK=off atlas migrate status --env local
```

### 4. Cập nhật migration hash

```bash
# Sau khi thêm migration mới
make migrate-hash
```

## Workflow

1. **Thay đổi model**: Sửa đổi struct trong `internal/domain/models.go`
2. **Sinh migration**: Chạy `make migrate-new NAME=<tên_migration>`
3. **Review migration**: Kiểm tra file SQL được sinh ra trong `migrations/`
4. **Apply migration**: Chạy `make migrate-up`

## Ví dụ thực tế

### Thêm một field mới vào Storage model:

```go
type Storage struct {
    // ... existing fields ...
    NewField string `gorm:"size:255" json:"new_field"`
}
```

Sau đó chạy:
```bash
make migrate-gen NAME=add_new_field_to_storage
```

Atlas sẽ tự động sinh file migration với SQL để thêm cột `new_field` vào bảng `storages`.

### Xóa một field:

```go
type Storage struct {
    // ... other fields ...
    // Xóa field OldField
}
```

Chạy:
```bash
make migrate-gen NAME=remove_old_field_from_storage
```

### Thay đổi kiểu dữ liệu:

```go
type Storage struct {
    // Thay đổi từ string sang int
    RateLimitRPM int64 `json:"rate_limit_rpm"` // Thay vì int
}
```

Chạy:
```bash
make migrate-gen NAME=change_rate_limit_type
```

## Các lệnh Makefile có sẵn

- `make migrate-gen NAME=<name>`: Sinh migration mới
- `make migrate-apply`: Apply tất cả migrations
- `make migrate-status`: Kiểm tra trạng thái migrations
- `make migrate-hash`: Cập nhật migration hash
- `make migrate-reset`: Reset database (cẩn thận!)

## Lưu ý

- ✅ **Luôn review file migration** trước khi apply
- ✅ **Backup database** trước khi apply migration trong production
- ✅ **Sử dụng `GOWORK=off`** để tránh lỗi workspace mode
- ✅ **Atlas sẽ tự động detect changes** và sinh SQL tương ứng
- ⚠️ **Không xóa migration files** đã được apply
- ⚠️ **Test migrations** trên development environment trước
