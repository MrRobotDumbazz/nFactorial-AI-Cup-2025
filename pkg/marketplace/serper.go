package marketplace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
)

const serperEndpoint = "https://google.serper.dev/search"

type SerperService struct {
	apiKey string
	client *http.Client
}

type serperRequest struct {
	Q string `json:"q"`
}

type serperResponse struct {
	OrganicResults []struct {
		Title       string `json:"title"`
		Link        string `json:"link"`
		Snippet     string `json:"snippet"`
		ImageURL    string `json:"imageUrl,omitempty"`
		Price       string `json:"price,omitempty"`
		Currency    string `json:"currency,omitempty"`
		Marketplace string `json:"marketplace,omitempty"`
	} `json:"organic"`
}

func NewSerperService() (*SerperService, error) {
	apiKey := os.Getenv("SERPER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SERPER_API_KEY environment variable is not set")
	}

	return &SerperService{
		apiKey: apiKey,
		client: &http.Client{},
	}, nil
}

func (s *SerperService) SearchProducts(ctx context.Context, query string) ([]types.Product, error) {
	// Формируем запрос к Serper API
	reqBody := serperRequest{
		Q: query,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", serperEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-KEY", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("serper API error: %s", string(body))
	}

	var serperResp serperResponse
	if err := json.Unmarshal(body, &serperResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Преобразуем результаты в нашу структуру Product
	var products []types.Product
	for _, item := range serperResp.OrganicResults {
		// Извлекаем цену из строки
		var price float64
		if item.Price != "" {
			fmt.Sscanf(strings.TrimSpace(strings.ReplaceAll(item.Price, ",", "")), "%f", &price)
		}

		// Определяем магазин из URL
		store := "unknown"
		if strings.Contains(item.Link, "kaspi.kz") {
			store = "kaspi"
		} else if strings.Contains(item.Link, "wildberries") {
			store = "wildberries"
		} else if strings.Contains(item.Link, "aliexpress") {
			store = "aliexpress"
		} else if strings.Contains(item.Link, "ozon") {
			store = "ozon"
		}

		product := types.Product{
			ID:          item.Link, // Используем URL как ID
			Title:       item.Title,
			Description: item.Snippet,
			Price:       price,
			URL:         item.Link,
			ImageURL:    item.ImageURL,
			Store:       store,
			Rating:      0, // У Serper нет информации о рейтинге
		}

		products = append(products, product)
	}

	return products, nil
}
