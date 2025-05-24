package main

import (
	"context"
	"log"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/marketplace"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var productService *marketplace.ProductService

func init() {
	// Инициализация AWS клиентов при холодном старте
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	productService = marketplace.NewProductService(dynamoClient)
}

func handleRequest(ctx context.Context, request types.ProductSearchRequest) (types.APIResponse, error) {
	// Проверка входных данных
	if len(request.Categories) == 0 {
		return types.APIResponse{
			Success: false,
			Error:   "at least one category is required",
		}, nil
	}

	// Поиск товаров
	products, err := productService.SearchProducts(ctx, request.Categories, request.PriceRange, request.Marketplace)
	if err != nil {
		log.Printf("Failed to search products: %v", err)
		return types.APIResponse{
			Success: false,
			Error:   "failed to search products",
		}, nil
	}

	return types.APIResponse{
		Success: true,
		Data: types.ProductSearchResponse{
			Products: products,
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
