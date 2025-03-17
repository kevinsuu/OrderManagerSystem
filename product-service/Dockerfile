# 使用多階段構建
FROM golang:1.22-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 複製整個專案
COPY . .

# 下載依賴
RUN go mod download

# 編譯
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# 最終階段
FROM alpine:latest

WORKDIR /app

# 從 builder 階段複製編譯好的二進制文件
COPY --from=builder /app/main .

# 設置容器啟動命令
CMD ["./main"]