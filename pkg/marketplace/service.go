package marketplace

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dyntypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ProductService struct {
	dynamoClient  *dynamodb.Client
	serperService *SerperService
	tableName     string
}

func NewProductService(dynamoClient *dynamodb.Client) *ProductService {
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		tableName = "products" // значение по умолчанию
	}

	serperService, err := NewSerperService()
	if err != nil {
		// Логируем ошибку, но продолжаем работу только с DynamoDB
		fmt.Printf("Failed to initialize Serper service: %v\n", err)
	}

	return &ProductService{
		dynamoClient:  dynamoClient,
		serperService: serperService,
		tableName:     tableName,
	}
}

func (s *ProductService) SearchProducts(ctx context.Context, categories []string, priceRange types.Range, marketplace string) ([]types.Product, error) {
	var allProducts []types.Product

	// Поиск в DynamoDB
	dbProducts, err := s.searchInDynamoDB(ctx, categories, priceRange, marketplace)
	if err != nil {
		fmt.Printf("DynamoDB search error: %v\n", err)
		// Продолжаем работу, даже если DynamoDB недоступен
	} else {
		allProducts = append(allProducts, dbProducts...)
	}

	// Если Serper сервис доступен, используем его для поиска
	if s.serperService != nil {
		// Формируем поисковый запрос
		query := s.buildSearchQuery(categories, priceRange, marketplace)

		serperProducts, err := s.serperService.SearchProducts(ctx, query)
		if err != nil {
			fmt.Printf("Serper search error: %v\n", err)
		} else {
			// Фильтруем результаты по категориям и ценовому диапазону
			for _, product := range serperProducts {
				if s.matchesFilters(product, categories, priceRange, marketplace) {
					allProducts = append(allProducts, product)
				}
			}
		}
	}

	// Удаляем дубликаты
	return s.removeDuplicates(allProducts), nil
}

func (s *ProductService) buildSearchQuery(categories []string, priceRange types.Range, marketplace string) string {
	// Базовый запрос
	query := strings.Join(categories, " OR ")

	// Добавляем ценовой диапазон
	if priceRange.Min > 0 || priceRange.Max > 0 {
		query += fmt.Sprintf(" price:%d..%d", int(priceRange.Min), int(priceRange.Max))
	}

	// Добавляем маркетплейс
	if marketplace != "" {
		query += " site:" + marketplace
	}

	return query
}

func (s *ProductService) matchesFilters(product types.Product, categories []string, priceRange types.Range, marketplace string) bool {
	// Проверка категории
	categoryMatch := false
	for _, category := range categories {
		if strings.Contains(strings.ToLower(product.Category), strings.ToLower(category)) {
			categoryMatch = true
			break
		}
	}
	if !categoryMatch {
		return false
	}

	// Проверка цены
	if priceRange.Min > 0 && product.Price < priceRange.Min {
		return false
	}
	if priceRange.Max > 0 && product.Price > priceRange.Max {
		return false
	}

	// Проверка маркетплейса
	if marketplace != "" && product.Store != marketplace {
		return false
	}

	return true
}

func (s *ProductService) removeDuplicates(products []types.Product) []types.Product {
	seen := make(map[string]bool)
	unique := make([]types.Product, 0)

	for _, product := range products {
		if !seen[product.ID] {
			seen[product.ID] = true
			unique = append(unique, product)
		}
	}

	return unique
}

func (s *ProductService) searchInDynamoDB(ctx context.Context, categories []string, priceRange types.Range, marketplace string) ([]types.Product, error) {
	// Создаем условия фильтрации
	filterExpr := "category IN (:categories)"
	if marketplace != "" {
		filterExpr += " AND marketplace = :marketplace"
	}
	filterExpr += " AND price BETWEEN :min_price AND :max_price"

	// Подготавливаем значения для условий
	categoriesValues := make([]dyntypes.AttributeValue, len(categories))
	for i, cat := range categories {
		categoriesValues[i] = &dyntypes.AttributeValueMemberS{Value: cat}
	}

	exprValues := map[string]dyntypes.AttributeValue{
		":categories": &dyntypes.AttributeValueMemberL{Value: categoriesValues},
		":min_price":  &dyntypes.AttributeValueMemberN{Value: fmt.Sprintf("%f", priceRange.Min)},
		":max_price":  &dyntypes.AttributeValueMemberN{Value: fmt.Sprintf("%f", priceRange.Max)},
	}

	if marketplace != "" {
		exprValues[":marketplace"] = &dyntypes.AttributeValueMemberS{Value: marketplace}
	}

	// Выполняем запрос к DynamoDB
	input := &dynamodb.ScanInput{
		TableName:                 &s.tableName,
		FilterExpression:          &filterExpr,
		ExpressionAttributeValues: exprValues,
	}

	result, err := s.dynamoClient.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan DynamoDB: %w", err)
	}

	// Преобразуем результаты в структуры Product
	var products []types.Product
	err = attributevalue.UnmarshalListOfMaps(result.Items, &products)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal products: %w", err)
	}

	// Помечаем продукты как из DynamoDB
	for i := range products {
		products[i].Store = "dynamodb"
	}

	return products, nil
}
