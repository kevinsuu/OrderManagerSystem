# 構建階段
FROM golang:1.21-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# 運行階段
FROM alpine:latest

# 安裝 CA 證書
RUN apk --no-cache add ca-certificates

# 設置工作目錄
WORKDIR /app

# 從構建階段複製二進制文件
COPY --from=builder /app/main .

# 創建存放 Firebase 憑證的目錄
RUN mkdir -p /app/credentials

# 暴露端口
EXPOSE 8082

# 設置環境變量
ENV SERVER_ADDRESS=:8082 \
    FIREBASE_CREDENTIALS=/app/credentials/firebase-credentials.json

# 運行應用
CMD ["./main"]
