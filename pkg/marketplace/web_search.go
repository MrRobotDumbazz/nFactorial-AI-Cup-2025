package marketplace

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
)

type WebSearchService struct {
	kaspiToken string
	aliToken   string
	wildToken  string
	ozonToken  string
	httpClient *http.Client
}

func NewWebSearchService(kaspiToken, aliToken, wildToken, ozonToken string) *WebSearchService {
	return &WebSearchService{
		kaspiToken: kaspiToken,
		aliToken:   aliToken,
		wildToken:  wildToken,
		ozonToken:  ozonToken,
		httpClient: &http.Client{},
	}
}

func (s *WebSearchService) SearchProducts(ctx context.Context, categories []string, priceRange types.Range, marketplace string) ([]types.Product, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allProducts []types.Product
	errors := make(chan error, 4)

	// Если указан конкретный маркетплейс, ищем только в нем
	if marketplace != "" {
		switch marketplace {
		case "kaspi":
			return s.searchKaspi(ctx, categories, priceRange)
		case "aliexpress":
			return s.searchAliExpress(ctx, categories, priceRange)
		case "wildberries":
			return s.searchWildberries(ctx, categories, priceRange)
		case "ozon":
			return s.searchOzon(ctx, categories, priceRange)
		default:
			return nil, fmt.Errorf("unsupported marketplace: %s", marketplace)
		}
	}

	// Если маркетплейс не указан, ищем везде параллельно
	if s.kaspiToken != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			products, err := s.searchKaspi(ctx, categories, priceRange)
			if err != nil {
				errors <- fmt.Errorf("kaspi search error: %v", err)
				return
			}
			mu.Lock()
			allProducts = append(allProducts, products...)
			mu.Unlock()
		}()
	}

	if s.aliToken != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			products, err := s.searchAliExpress(ctx, categories, priceRange)
			if err != nil {
				errors <- fmt.Errorf("aliexpress search error: %v", err)
				return
			}
			mu.Lock()
			allProducts = append(allProducts, products...)
			mu.Unlock()
		}()
	}

	if s.wildToken != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			products, err := s.searchWildberries(ctx, categories, priceRange)
			if err != nil {
				errors <- fmt.Errorf("wildberries search error: %v", err)
				return
			}
			mu.Lock()
			allProducts = append(allProducts, products...)
			mu.Unlock()
		}()
	}

	if s.ozonToken != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			products, err := s.searchOzon(ctx, categories, priceRange)
			if err != nil {
				errors <- fmt.Errorf("ozon search error: %v", err)
				return
			}
			mu.Lock()
			allProducts = append(allProducts, products...)
			mu.Unlock()
		}()
	}

	// Ждем завершения всех поисков
	wg.Wait()
	close(errors)

	// Проверяем ошибки
	var errStrings []string
	for err := range errors {
		errStrings = append(errStrings, err.Error())
	}
	if len(errStrings) > 0 {
		log.Printf("Some marketplace searches failed: %v", strings.Join(errStrings, "; "))
	}

	return allProducts, nil
}

func (s *WebSearchService) searchKaspi(ctx context.Context, categories []string, priceRange types.Range) ([]types.Product, error) {
	// Преобразуем категории в формат Kaspi
	kaspiCategories := make([]string, 0)
	for _, cat := range categories {
		if kaspiCat, ok := types.CategoryMappings["kaspi"][cat]; ok {
			kaspiCategories = append(kaspiCategories, kaspiCat)
		}
	}

	// Формируем URL для API Kaspi
	baseURL := "https://kaspi.kz/shop/api/products/search"
	params := url.Values{}
	params.Add("categories", strings.Join(kaspiCategories, ","))
	params.Add("price_from", fmt.Sprintf("%.0f", priceRange.Min))
	params.Add("price_to", fmt.Sprintf("%.0f", priceRange.Max))

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.kaspiToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []types.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}

	// Добавляем маркетплейс к каждому продукту
	for i := range products {
		products[i].Store = "kaspi"
	}

	return products, nil
}

func (s *WebSearchService) searchAliExpress(ctx context.Context, categories []string, priceRange types.Range) ([]types.Product, error) {
	// Аналогичная реализация для AliExpress
	// Используем AliExpress Open API
	aliCategories := make([]string, 0)
	for _, cat := range categories {
		if aliCat, ok := types.CategoryMappings["aliexpress"][cat]; ok {
			aliCategories = append(aliCategories, aliCat)
		}
	}

	baseURL := "https://api.aliexpress.com/v2/products/search"
	params := url.Values{}
	params.Add("categories", strings.Join(aliCategories, ","))
	params.Add("price_min", fmt.Sprintf("%.2f", priceRange.Min))
	params.Add("price_max", fmt.Sprintf("%.2f", priceRange.Max))

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.aliToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []types.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}

	for i := range products {
		products[i].Store = "aliexpress"
	}

	return products, nil
}

func (s *WebSearchService) searchWildberries(ctx context.Context, categories []string, priceRange types.Range) ([]types.Product, error) {
	// Реализация для Wildberries
	wbCategories := make([]string, 0)
	for _, cat := range categories {
		if wbCat, ok := types.CategoryMappings["wildberries"][cat]; ok {
			wbCategories = append(wbCategories, wbCat)
		}
	}

	baseURL := "https://suppliers-api.wildberries.ru/api/v3/products"
	params := url.Values{}
	params.Add("subjects", strings.Join(wbCategories, ","))
	params.Add("priceFrom", fmt.Sprintf("%.0f", priceRange.Min))
	params.Add("priceTo", fmt.Sprintf("%.0f", priceRange.Max))

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", s.wildToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []types.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}

	for i := range products {
		products[i].Store = "wildberries"
	}

	return products, nil
}

func (s *WebSearchService) searchOzon(ctx context.Context, categories []string, priceRange types.Range) ([]types.Product, error) {
	// Реализация для Ozon
	ozonCategories := make([]string, 0)
	for _, cat := range categories {
		if ozonCat, ok := types.CategoryMappings["ozon"][cat]; ok {
			ozonCategories = append(ozonCategories, ozonCat)
		}
	}

	baseURL := "https://api-seller.ozon.ru/v3/product/list"

	// Формируем тело запроса
	requestBody := map[string]interface{}{
		"categories": ozonCategories,
		"price": map[string]float64{
			"from": priceRange.Min,
			"to":   priceRange.Max,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-Id", s.ozonToken)
	req.Header.Set("Api-Key", s.ozonToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []types.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}

	for i := range products {
		products[i].Store = "ozon"
	}

	return products, nil
}
