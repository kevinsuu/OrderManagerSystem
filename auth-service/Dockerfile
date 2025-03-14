FROM golang:1.22-alpine AS builder

WORKDIR /app

# 首先只複製 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 再複製其餘源代碼
COPY . .

# 構建
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service ./cmd/main.go

# 最終階段
FROM alpine:latest

WORKDIR /app

# 從構建階段複製二進制文件
COPY --from=builder /app/auth-service .

# 複製 Firebase 憑證
COPY order-manager-system-a6931-firebase-adminsdk-fbsvc-65d3904bc6.json /app/firebase-credentials.json

# 暴露端口
EXPOSE 8081

# 運行應用
CMD ["./auth-service"]
