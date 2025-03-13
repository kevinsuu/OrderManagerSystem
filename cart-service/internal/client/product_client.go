package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	Stock       int            `json:"stock"` // 修改自 StockCount
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

func NewProductClient(baseURL string) ProductClient {
	return &productClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *productClient) GetProduct(ctx context.Context, productID string) (*ProductInfo, error) {
	url := fmt.Sprintf("%s/api/v1/products/%s", c.baseURL, productID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 直接使用完整的 Authorization header
	if token := ctx.Value(TokenKey); token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("%v", token))
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

	// 處理不同的狀態碼
	switch resp.StatusCode {
	case http.StatusOK:
		var product ProductInfo
		if err := json.Unmarshal(body, &product); err != nil {
			return nil, fmt.Errorf("unmarshal response failed: %w", err)
		}
		return &product, nil
	case http.StatusNotFound:
		return nil, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized: %s", string(body))
	default:
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
}
