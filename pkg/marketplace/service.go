package marketplace

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dyntypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ProductService struct {
	dynamoClient *dynamodb.Client
	aiSearch     *AISearchService
	tableName    string
}

func NewProductService(dynamoClient *dynamodb.Client) *ProductService {
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		tableName = "products" // значение по умолчанию
	}

	// Получаем токен Serper API из переменных окружения
	serperToken := os.Getenv("SERPER_API_TOKEN")
	aiSearch := NewAISearchService(serperToken)

	return &ProductService{
		dynamoClient: dynamoClient,
		aiSearch:     aiSearch,
		tableName:    tableName,
	}
}

func (s *ProductService) SearchProducts(ctx context.Context, categories []string, priceRange types.Range, marketplace string) ([]types.Product, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allProducts []types.Product
	errors := make(chan error, 2)

	// Поиск в DynamoDB
	wg.Add(1)
	go func() {
		defer wg.Done()
		products, err := s.searchInDynamoDB(ctx, categories, priceRange, marketplace)
		if err != nil {
			errors <- fmt.Errorf("dynamodb search error: %v", err)
			return
		}
		mu.Lock()
		allProducts = append(allProducts, products...)
		mu.Unlock()
	}()

	// Поиск через AI Search
	wg.Add(1)
	go func() {
		defer wg.Done()
		products, err := s.aiSearch.SearchProducts(ctx, categories, priceRange, marketplace)
		if err != nil {
			errors <- fmt.Errorf("ai search error: %v", err)
			return
		}
		mu.Lock()
		allProducts = append(allProducts, products...)
		mu.Unlock()
	}()

	// Ждем завершения всех поисков
	wg.Wait()
	close(errors)

	// Проверяем ошибки
	var errStrings []string
	for err := range errors {
		log.Printf("Search error: %v", err)
		errStrings = append(errStrings, err.Error())
	}

	// Удаляем дубликаты по ID продукта
	uniqueProducts := make(map[string]types.Product)
	for _, product := range allProducts {
		// Предпочитаем результаты из AI поиска (они более актуальные)
		if _, exists := uniqueProducts[product.ID]; !exists || product.Store != "dynamodb" {
			uniqueProducts[product.ID] = product
		}
	}

	// Преобразуем обратно в слайс
	result := make([]types.Product, 0, len(uniqueProducts))
	for _, product := range uniqueProducts {
		result = append(result, product)
	}

	return result, nil
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
