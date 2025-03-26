package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ContextKey 自定義 context key 類型
type ContextKey string

const (
	// TokenKey 用於存儲 token 的 context key
	TokenKey ContextKey = "token"
)

type ProductClient interface {
	GetProduct(ctx context.Context, productID string) (*ProductInfo, error)
	GetProductImageAsBase64(ctx context.Context, imageURL string) (string, error)
}

type ProductImage struct {
	ID        string `json:"id"`
	ProductID string `json:"productId"`
	URL       string `json:"url"`
	Sort      int    `json:"sort"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ProductInfo struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Stock       int            `json:"stock"`
	Status      string         `json:"status"`
	CategoryID  string         `json:"categoryId"`
	Images      []ProductImage `json:"images"`
	Attributes  []interface{}  `json:"attributes"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
}

type productClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewProductClient 創建新的 ProductClient 實例
func NewProductClient(baseURL string) ProductClient {
	// 如果沒有提供 baseURL，使用默認值
	if baseURL == "" {
		baseURL = "https://ordermanagersystem-product-service.onrender.com"
		fmt.Printf("No PRODUCT_SERVICE_URL provided, using default: %s", baseURL)
	}

	return &productClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *productClient) GetProduct(ctx context.Context, productID string) (*ProductInfo, error) {
	url := fmt.Sprintf("%s/api/v1/products/%s", c.baseURL, productID)

	// 添加日誌來查看實際的 URL
	log.Printf("Requesting product from URL: %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 添加日誌來查看 token
	if token := ctx.Value(TokenKey); token != nil {
		log.Printf("Using token: %v", token)
		req.Header.Set("Authorization", fmt.Sprintf("%v", token))
	} else {
		log.Printf("No token found in context")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	log.Printf("Product service response status: %d", resp.StatusCode)
	log.Printf("Product service response body: %s", string(body))

	// 處理不同的狀態碼
	switch resp.StatusCode {
	case http.StatusOK:
		// 定義包裝響應結構體
		type ProductResponse struct {
			Data    ProductInfo `json:"data"`
			Message string      `json:"message"`
			Success bool        `json:"success"`
		}

		var response ProductResponse
		if err := json.Unmarshal(body, &response); err != nil {
			// 嘗試直接解析為ProductInfo
			var product ProductInfo
			if err2 := json.Unmarshal(body, &product); err2 != nil {
				return nil, fmt.Errorf("unmarshal response failed: %w, tried direct unmarshal: %w", err, err2)
			}
			return &product, nil
		}

		// 日誌輸出解析後的產品信息
		log.Printf("Parsed product info: %+v", response.Data)

		return &response.Data, nil
	case http.StatusNotFound:
		return nil, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized: %s", string(body))
	default:
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
}

// 新增一個方法來獲取產品圖片並轉換為base64
func (c *productClient) GetProductImageAsBase64(ctx context.Context, imageURL string) (string, error) {
	if imageURL == "" {
		return "", nil
	}

	// 使用HTTP GET請求獲取圖片
	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("create image request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("image request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get image, status: %d", resp.StatusCode)
	}

	// 讀取圖片數據
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read image data failed: %w", err)
	}

	// 轉換為base64
	base64Data := "data:" + resp.Header.Get("Content-Type") + ";base64," + base64.StdEncoding.EncodeToString(imgData)

	return base64Data, nil
}
