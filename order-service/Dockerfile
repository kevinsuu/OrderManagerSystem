FROM golang:1.22-alpine AS builder

# 安裝基本構建工具
RUN apk add --no-cache git

WORKDIR /app

# 設置 Go 環境變數以優化構建
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 複製並下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製源代碼
COPY . .

# 構建（添加編譯優化標誌）
RUN go build -ldflags="-w -s" -o /app/order-service ./cmd/main.go

# 最終階段
FROM alpine:latest

# 安裝 CA 證書
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 從構建階段複製二進制文件
COPY --from=builder /app/order-service .

# 暴露端口
EXPOSE 8083

# 運行應用
CMD ["./order-service"]
