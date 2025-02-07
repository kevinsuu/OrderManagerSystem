FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安裝必要的系統依賴
RUN apk add --no-cache gcc musl-dev

# 複製 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/notification-service ./cmd/main.go

# 最終階段
FROM alpine:latest

WORKDIR /app

# 安裝必要的運行時依賴
RUN apk add --no-cache ca-certificates tzdata

# 設置時區
ENV TZ=Asia/Taipei

# 從構建階段複製二進制文件
COPY --from=builder /app/notification-service .

# 設置默認環境變量
ENV SERVER_ADDRESS=:8084 \
    DB_HOST=postgres \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=notification_service \
    DB_SSLMODE=disable \
    REDIS_ADDRESS=redis:6379 \
    REDIS_PASSWORD="" \
    REDIS_DB=0 \
    EMAIL_HOST=smtp.gmail.com \
    EMAIL_PORT=587 \
    SMS_PROVIDER=twilio \
    PUSH_PROVIDER=firebase

# 暴露端口
EXPOSE 8084

# 健康檢查
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8084/health || exit 1

# 運行應用
CMD ["./notification-service"]
