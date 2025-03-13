FROM golang:1.23-alpine AS builder

WORKDIR /app

# 複製所有源代碼
COPY . .

# 下載依賴並構建
RUN go mod download -x && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/cart-service ./cmd/main.go

# 最終階段
FROM alpine:latest

WORKDIR /app

# 從構建階段複製二進制文件
COPY --from=builder /app/cart-service .

# 暴露端口
EXPOSE 8086

# 運行應用
CMD ["./cart-service"]
