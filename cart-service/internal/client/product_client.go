package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProductClient interface {
	GetProduct(ctx context.Context, productID string) (*ProductInfo, error)
}

type ProductInfo struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	Price      float64 `json:"price"`
	StockCount int     `json:"stockCount"`
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var product ProductInfo
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &product, nil
}
