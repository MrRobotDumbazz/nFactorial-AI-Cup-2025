package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/marketplace"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
	"github.com/aws/aws-lambda-go/events"
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

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
		"Content-Type":                "application/json",
	}

	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
		}, nil
	}

	var searchRequest types.ProductSearchRequestApi
	if err := json.Unmarshal([]byte(request.Body), &searchRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"Invalid request body"}`,
			Headers:    headers,
		}, nil
	}

	if len(searchRequest.Categories) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"at least one category is required"}`,
			Headers:    headers,
		}, nil
	}

	var priceRange types.Range
	if searchRequest.PriceRange != nil {
		priceRange = *searchRequest.PriceRange
	}

	products, err := productService.SearchProducts(ctx, searchRequest.Categories, priceRange, searchRequest.Marketplace)
	if err != nil {
		log.Printf("Failed to search products: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"failed to search products"}`,
			Headers:    headers,
		}, nil
	}

	response := types.ApiResponse{
		Success: true,
		Data: types.ProductSearchResponseApi{
			Products: products,
		},
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"Failed to marshal response"}`,
			Headers:    headers,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
