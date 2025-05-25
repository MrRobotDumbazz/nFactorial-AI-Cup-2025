package marketplace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
)

type AISearchService struct {
	serperToken string
	httpClient  *http.Client
}

type SerperRequest struct {
	Q          string `json:"q"`
	GL         string `json:"gl"`   // Геолокация (например, "kz" для Казахстана)
	Num        int    `json:"num"`  // Количество результатов
	SearchType string `json:"type"` // shopping для поиска товаров
}

type SerperResponse struct {
	Shopping []struct {
		Title    string `json:"title"`
		Link     string `json:"link"`
		Price    string `json:"price"`
		Source   string `json:"source"`
		ImageURL string `json:"imageUrl"`
	} `json:"shopping"`
}

func NewAISearchService(serperToken string) *AISearchService {
	return &AISearchService{
		serperToken: serperToken,
		httpClient:  &http.Client{},
	}
}

func (s *AISearchService) SearchProducts(ctx context.Context, categories []string, priceRange types.Range, marketplace string) ([]types.Product, error) {
	// Формируем поисковый запрос
	query := s.buildSearchQuery(categories, priceRange, marketplace)

	// Создаем запрос к Serper API
	reqBody := SerperRequest{
		Q:          query,
		GL:         "kz",       // Ищем в Казахстане
		Num:        20,         // Получаем 20 результатов
		SearchType: "shopping", // Поиск по товарам
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://google.serper.dev/shopping", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", s.serperToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var serperResp SerperResponse
	if err := json.Unmarshal(body, &serperResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Преобразуем результаты в наш формат
	return s.convertToProducts(serperResp.Shopping, categories[0])
}

func (s *AISearchService) buildSearchQuery(categories []string, priceRange types.Range, marketplace string) string {
	query := strings.Join(categories, " ")

	// Добавляем ценовой диапазон
	if priceRange.Min > 0 || priceRange.Max > 0 {
		query += fmt.Sprintf(" price %d-%d тенге", int(priceRange.Min), int(priceRange.Max))
	}

	// Добавляем маркетплейс если указан
	if marketplace != "" {
		query += " site:" + marketplace
	}

	// Добавляем "купить в Казахстане" для релевантности
	query += " купить в Казахстане"

	return query
}

func (s *AISearchService) convertToProducts(results []struct {
	Title    string `json:"title"`
	Link     string `json:"link"`
	Price    string `json:"price"`
	Source   string `json:"source"`
	ImageURL string `json:"imageUrl"`
}, category string) ([]types.Product, error) {
	var products []types.Product

	for _, result := range results {
		price, _ := s.extractPrice(result.Price)

		product := types.Product{
			ID:          generateProductID(result.Link),
			Title:       result.Title,
			Description: result.Title, // Используем заголовок как описание
			Price:       price,
			URL:         result.Link,
			ImageURL:    result.ImageURL,
			Store:       result.Source,
			Category:    category,
			Rating:      0, // У нас нет рейтинга из поиска
		}

		products = append(products, product)
	}

	return products, nil
}

func (s *AISearchService) extractPrice(priceStr string) (float64, error) {
	// Удаляем все символы кроме цифр и точки
	re := regexp.MustCompile(`[^\d.]`)
	cleanPrice := re.ReplaceAllString(priceStr, "")

	// Конвертируем в float64
	price, err := strconv.ParseFloat(cleanPrice, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func generateProductID(url string) string {
	// Создаем ID из URL, убирая протокол и домен
	parts := strings.Split(url, "/")
	if len(parts) > 3 {
		return strings.Join(parts[3:], "-")
	}
	return strings.Join(parts, "-")
}
