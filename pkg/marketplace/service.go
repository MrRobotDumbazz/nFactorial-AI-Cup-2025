package marketplace

import (
	"context"
	"fmt"
	"os"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dyntypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ProductService struct {
	dynamoClient *dynamodb.Client
	tableName    string
}

func NewProductService(dynamoClient *dynamodb.Client) *ProductService {
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		tableName = "products" // значение по умолчанию
	}
	return &ProductService{
		dynamoClient: dynamoClient,
		tableName:    tableName,
	}
}

func (s *ProductService) SearchProducts(ctx context.Context, categories []string, priceRange types.Range, marketplace string) ([]types.Product, error) {
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

	return products, nil
}
