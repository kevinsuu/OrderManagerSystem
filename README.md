# 訂單管理系統 (Order Management System)

這是一個基於微服務架構的訂單管理系統，使用現代化的技術棧和最佳實踐來構建。

## 技術清單

- **後端框架**: Go (Gin/gRPC)
- **資料庫**: PostgreSQL
- **快取**: Redis
- **容器化**: Docker
- **容器編排**: Kubernetes
- **認證**: JWT
- **反向代理**: Nginx
- **CI/CD**: GitHub Actions
- **雲端部署**: Render

## 系統架構

系統包含以下微服務：

1. **API Gateway**: 
   - 統一的 API 入口
   - 請求路由和負載均衡
   - 請求限流和認證

2. **Auth Service**:
   - 用戶認證和授權
   - JWT token 管理
   - 用戶會話管理

3. **Order Service**:
   - 訂單創建和管理
   - 訂單狀態追蹤
   - 訂單歷史記錄

4. **Product Service**:
   - 產品目錄管理
   - 庫存管理
   - 價格管理

5. **Payment Service**:
   - 支付處理
   - 退款處理
   - 交易記錄

6. **Notification Service**:
   - 郵件通知
   - 系統通知
   - 事件推送

## 開發環境設置

### 前置需求

- Go 1.21+
- Docker
- Kubernetes
- PostgreSQL
- Redis

### 本地開發

1. 克隆專案
```bash
git clone https://github.com/kevinsuu/OrderManagerSystem.git
cd OrderManagerSystem
```

2. 安裝依賴
```bash
make deps
```

3. 啟動服務
```bash
make run
```

### Docker 部署

```bash
make docker-build
make docker-run
```

### Kubernetes 部署

```bash
make k8s-deploy
```

## API 文檔

API 文檔使用 Swagger 生成，可在本地開發環境訪問：
- REST API: http://localhost:8080/swagger/index.html
- gRPC: 參考 proto 文件

## 監控和日誌

- Prometheus 用於指標收集
- Grafana 用於監控面板
- ELK Stack 用於日誌管理

## 貢獻指南

1. Fork 專案
2. 創建特性分支
3. 提交變更
4. 推送到分支
5. 創建 Pull Request

## 授權

MIT License
