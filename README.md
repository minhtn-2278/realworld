# Realworld API

REST API viết bằng Go, Echo v5, GORM v2, PostgreSQL và Redis. API được tổ chức theo luồng `handler → service → repository` và có Swagger UI để khám phá endpoint.

## Yêu cầu

- Go 1.26+
- Docker và Docker Compose
- PostgreSQL và Redis, hoặc dùng các service trong `docker-compose.yaml`
- `psql` nếu muốn chạy seed data bằng lệnh SQL

## Cấu trúc thư mục

```text
cmd/
├── app/                 # HTTP application entry point
└── migrate/             # CLI chạy migration up/down
config/                  # Load biến môi trường và cấu hình ứng dụng
docs/                    # Swagger spec được sinh tự động
internal/
├── cache/               # Redis cache store
├── dto/                 # Request/response DTO và mapper response
├── handlers/            # HTTP handlers, routes và HTTP error handler
├── middleware/          # JWT auth và resource-owner middleware
├── models/              # GORM model và quan hệ dữ liệu
├── repositories/        # Persistence/query bằng GORM
├── services/            # Business logic/use cases
└── utils/               # Error, validation và pagination helpers
migrations/              # GORMigrate migrations
pkg/db/                  # PostgreSQL/GORM connection setup
seed.sql                 # Dữ liệu mẫu dùng cho local development
docker-compose.yaml      # PostgreSQL, Redis và Adminer cho development
```

## Cài đặt và chạy local

1. Tạo file cấu hình từ template:

   ```bash
   cp .env.example .env
   ```

2. Chỉnh `JWT_SECRET` trong `.env` thành chuỗi bí mật riêng. Các giá trị PostgreSQL/Redis mặc định tương thích với Docker Compose.

3. Khởi động PostgreSQL và Redis:

   ```bash
   docker compose up -d realworld-postgres realworld-redis
   ```

   Adminer là tùy chọn, có thể khởi động thêm với:

   ```bash
   docker compose up -d realworld-adminer
   ```

4. Chạy migrations:

   ```bash
   go run ./cmd/migrate -direction=up
   ```

5. Nạp dữ liệu mẫu (tùy chọn):

   ```bash
   psql "$DATABASE_URL" -f seed.sql
   ```

   Nếu PostgreSQL đang chạy bằng Docker Compose, dùng:

   ```bash
   docker compose exec -T realworld-postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < seed.sql
   ```

   `seed.sql` có thể chạy nhiều lần. Các user mẫu là `seed_alice`, `seed_bob`, `seed_carol`; mật khẩu chung là `password`.

6. Khởi động HTTP API:

   ```bash
   go run ./cmd/app
   ```

Mặc định API lắng nghe tại `http://localhost:3001` theo `.env.example`.

## Swagger

Sau khi khởi động app, mở Swagger UI tại:

```text
http://localhost:3001/swagger/index.html
```

Swagger spec được lưu ở `docs/swagger.json` và `docs/swagger.yaml`. Sau khi thêm hoặc thay đổi annotation API, sinh lại spec bằng:

```bash
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init --parseInternal -g cmd/app/main.go -o docs
```

## Các lệnh thường dùng

| Mục đích | Lệnh |
| --- | --- |
| Chạy migration mới | `go run ./cmd/migrate -direction=up` |
| Rollback migration gần nhất | `go run ./cmd/migrate -direction=down` |
| Chạy API | `go run ./cmd/app` |
| Format file Go đã sửa | `gofmt -w <files>` |
| Chạy toàn bộ test | `go test ./...` |
| Build toàn bộ project | `go build ./...` |
| Sinh Swagger spec | `go run github.com/swaggo/swag/cmd/swag@v1.16.6 init --parseInternal -g cmd/app/main.go -o docs` |
| Dừng development services | `docker compose down` |

Nếu Go build cache mặc định không có quyền ghi trong môi trường local/CI, dùng:

```bash
GOCACHE=/tmp/realworld-go-cache go test ./...
```

## Cấu hình môi trường

| Biến | Mô tả |
| --- | --- |
| `HTTP_ADDR` | Địa chỉ HTTP, ví dụ `:3001` |
| `DATABASE_URL` | PostgreSQL connection URL; bắt buộc |
| `REDIS_URL` | Redis connection URL |
| `JWT_SECRET` | Secret ký JWT; bắt buộc |
| `POSTGRES_DB` | Tên database dùng bởi Docker Compose |
| `POSTGRES_USER` | PostgreSQL user dùng bởi Docker Compose |
| `POSTGRES_PASSWORD` | PostgreSQL password dùng bởi Docker Compose |
| `POSTGRES_PORT` | PostgreSQL host port |
| `REDIS_PORT` | Redis host port |

Không commit file `.env` hoặc sử dụng secret development trong môi trường production.
