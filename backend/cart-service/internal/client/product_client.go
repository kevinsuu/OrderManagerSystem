package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kevinsuu/OrderManagerSystem/cart-service/internal/model"
)

// ContextKey 自定義 context key 類型
type ContextKey string

const (
	// TokenKey 用於存儲 token 的 context key
	TokenKey ContextKey = "token"
)

// ProductClient 提供與產品服務交互的功能
type ProductClient interface {
	GetProduct(ctx context.Context, productID string) (*ProductInfo, error)
	GetProductImageAsBase64(ctx context.Context, imageURL string) (string, error)
	GetProductById(ctx context.Context, productId string) (*model.ProductInfo, error)
}

type ProductImage struct {
	ID        string `json:"id"`
	ProductID string `json:"productId"`
	URL       string `json:"url"`
	Data      string `json:"data"` // 添加 Data 字段存儲 base64 圖片數據
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

// productClient 實現 ProductClient 接口
type productClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewProductClient 創建一個新的產品客戶端
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

// 添加一個函數來獲取並轉換圖片為 base64
func (c *productClient) GetProductImageAsBase64(ctx context.Context, imageURL string) (string, error) {
	if imageURL == "" {
		log.Printf("Empty image URL provided")
		return "", nil
	}

	log.Printf("Fetching image from URL: %s", imageURL)

	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		log.Printf("Failed to create request for image: %v", err)
		return "", fmt.Errorf("create image request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to fetch image: %v", err)
		return "", fmt.Errorf("image request failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Image fetch response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		log.Printf("Image fetch failed with status %d: %s", resp.StatusCode, string(respBody))
		return "", fmt.Errorf("failed to get image, status: %d", resp.StatusCode)
	}

	// 讀取圖片數據
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read image data: %v", err)
		return "", fmt.Errorf("read image data failed: %w", err)
	}

	log.Printf("Successfully read image data, size: %d bytes", len(imgData))

	// 確定 MIME 類型
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(imgData)
		log.Printf("Detected content type: %s", contentType)
	}

	// 轉換為 base64
	base64Data := "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(imgData)

	log.Printf("Successfully converted image to base64 (length: %d)", len(base64Data))
	return base64Data, nil
}

// GetProductById 通過 ID 獲取產品詳細資訊
func (c *productClient) GetProductById(ctx context.Context, productId string) (*model.ProductInfo, error) {
	url := fmt.Sprintf("%s/api/v1/products/%s", c.baseURL, productId)

	log.Printf("Requesting product details from URL: %s", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 添加授權頭
	if token := ctx.Value(TokenKey); token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("%v", token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 讀取並記錄響應
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Product service response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from product service: %s", resp.Status)
	}

	// 從響應中提取產品數據
	// 先解析整體響應結構
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 檢查success字段
	success, ok := response["success"].(bool)
	if !ok || !success {
		return nil, fmt.Errorf("product service returned unsuccessful response")
	}

	// 提取data部分
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no data field in response")
	}

	// 創建我們簡化的ProductInfo結構
	productInfo := &model.ProductInfo{
		ID:    productId,
		Name:  "",
		Price: 0,
	}

	// 嘗試獲取產品名稱
	if name, ok := data["name"].(string); ok {
		productInfo.Name = name
	}

	// 嘗試獲取產品價格
	if price, ok := data["price"].(float64); ok {
		productInfo.Price = price
	}

	// 嘗試獲取圖片列表
	if images, ok := data["images"].([]interface{}); ok && len(images) > 0 {
		// 將圖片URL轉換為字符串數組
		imageUrls := make([]string, 0, len(images))
		for _, img := range images {
			if imgMap, ok := img.(map[string]interface{}); ok {
				if url, ok := imgMap["url"].(string); ok && url != "" {
					imageUrls = append(imageUrls, url)
				}
			}
		}
		productInfo.Images = imageUrls
	}

	// 嘗試獲取其他可能有用的字段
	if createdAt, ok := data["createdAt"].(string); ok {
		productInfo.CreatedAt = createdAt
	}

	if updatedAt, ok := data["updatedAt"].(string); ok {
		productInfo.UpdatedAt = updatedAt
	}

	log.Printf("Successfully parsed product info: %+v", productInfo)
	return productInfo, nil
}
