FROM golang:1.23-alpine AS builder

WORKDIR /app

# 複製所有源代碼
COPY . .

# 下載依賴並構建
RUN go mod download -x && \
  go mod tidy && \
  CGO_ENABLED=0 GOOS=linux go build -o /app/notification-service ./cmd/main.go

# 最終階段
FROM alpine:latest

WORKDIR /app

# 從構建階段複製二進制文件
COPY --from=builder /app/notification-service .

# 暴露端口
EXPOSE 8084

# 運行應用
CMD ["./notification-service"]
