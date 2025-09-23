# Hướng dẫn phát triển API mới

Tài liệu này mô tả các bước chuẩn để thêm một API mới vào `chorus-controller` theo kiến trúc hiện tại.

## 1) Thiết kế API
- Xác định: phương thức (GET/POST/PUT/DELETE), endpoint, input, output, quy tắc auth.
- Quy ước header: sử dụng `Authorization: Token <TOKEN>`.
- Phân loại quyền:
  - Public: không cần token (ví dụ read-only công khai).
  - System-only: yêu cầu system token (quản trị).
  - Normal token: yêu cầu token thường (ghi dữ liệu).

## 2) Khai báo kiểu dữ liệu domain
- Thêm DTO/Models vào `internal/domain/` nếu cần:
  - Request/Response struct: đặt tên rõ ràng, có json tags.
  - Nếu thay đổi schema DB: cập nhật `models.go` và xem mục Migration bên dưới.

Ví dụ (request/response):
```go
// internal/domain/models.go (hoặc file riêng trong domain nếu chỉ là DTO)
type CreateThingRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type ThingResponse struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
}
```

## 3) Cập nhật interface service
- Thêm method vào `internal/domain/interfaces.go` trong interface phù hợp (ví dụ `StorageService`, `ReplicationService`, `TokenService`, ...).
- Giữ interface gọn rõ, có context.

```go
// internal/domain/interfaces.go
type ThingService interface {
    CreateThing(ctx context.Context, req *CreateThingRequest) (*ThingResponse, error)
}
```

## 4) Implement service
- Thêm logic vào `internal/service/` (tạo file mới nếu cần).
- Thực hiện validate nghiệp vụ, gọi repository.
- Trả về lỗi qua package `internal/errors` để handler hiển thị chuẩn.

```go
// internal/service/thing.go
type ThingService struct { repo *repository.ThingRepository }

func (s *ThingService) CreateThing(ctx context.Context, req *domain.CreateThingRequest) (*domain.ThingResponse, error) {
    if strings.TrimSpace(req.Name) == "" {
        return nil, errors.NewBadRequestError("name is required", nil)
    }
    // gọi repo, map model -> response
    // ...
    return resp, nil
}
```

## 5) Repository (nếu cần DB)
- Thêm repository vào `internal/repository/`, dùng GORM.
- Cung cấp các hàm CRUD cần thiết, nhận context.

```go
// internal/repository/thing_db.go
type ThingRepository struct{}
func (r *ThingRepository) Create(ctx context.Context, m *domain.Thing) error { /* db.DB().Create(m) */ return nil }
```

## 6) HTTP Handler
- Thêm handler vào `internal/handler/` hoặc mở rộng handler hiện có.
- Dùng `gin.Context`, bind/validate input, gọi service, trả JSON.
- Dùng `middleware.ErrorResponse/HandleError` cho lỗi.
- Thêm swagger comments đúng định dạng ở trên handler (Summary/Description/Tags/Params/Success/Failure/Router).

```go
// internal/handler/thing.go
// @Summary      Create new Thing
// @Description  ...
// @Tags         thing
// @Accept       json
// @Produce      json
// @Security     TokenAuth
// @Param        request  body  domain.CreateThingRequest  true  "payload"
// @Success      201      {object}  domain.ThingResponse
// @Failure      400      {object}  map[string]interface{}
// @Router       /things [post]
func (h *ThingHandler) Create(c *gin.Context) {
    var req domain.CreateThingRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, middleware.ErrorResponse(err))
        return
    }
    resp, err := h.svc.CreateThing(c.Request.Context(), &req)
    if err != nil {
        middleware.HandleError(c, err)
        return
    }
    c.JSON(http.StatusCreated, resp)
}
```

## 7) Routing & Middleware
- Đăng ký route trong `internal/server/server.go`.
- Chọn middleware phù hợp:
  - Public: không cần middleware.
  - System-only: `middleware.SystemTokenAuth(s.tokenService)`.
  - Normal: `middleware.AuthMiddleware(s.tokenService)`.

```go
// server.go (ví dụ)
systemProtected := r.Group("/")
systemProtected.Use(middleware.SystemTokenAuth(s.tokenService))
{
    systemProtected.POST("/things", thingHandler.Create)
}
```

## 8) Swagger Docs
- Sau khi thêm swagger comments, cập nhật tài liệu:
```bash
make swagger
```
- Mở `http://localhost:8081/swagger/index.html` để kiểm tra.

## 9) Migration (khi đổi schema)
- Khi sửa `internal/domain/models.go`, sinh migration:
```bash
make migrate-new NAME=<ten_migration>
make migrate-up
```
- Xem `docs/getting-started/migrations.md` để biết chi tiết.

## 10) Kiểm thử
- Unit test cho service/repository nếu có thể: `make test`.
- Thử API thủ công với curl:
```bash
curl -X POST http://localhost:8081/things \
  -H "Content-Type: application/json" \
  -H "Authorization: Token <TOKEN>" \
  -d '{"name":"demo","description":"..."}'
```

## 11) Checklist nhanh
- [ ] Cập nhật domain DTO/models
- [ ] Cập nhật interface service
- [ ] Implement service + (repository nếu cần)
- [ ] Viết handler + swagger comments
- [ ] Đăng ký route + middleware phù hợp
- [ ] Sinh/áp dụng migration (nếu đổi schema)
- [ ] Cập nhật Swagger (`make swagger`)
- [ ] Viết test cơ bản (`make test`)
- [ ] Thử nghiệm qua curl/Swagger UI
