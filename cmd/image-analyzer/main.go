package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/analyzer"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
)

var imageAnalyzer *analyzer.ImageAnalyzer

func init() {
	// Инициализация AWS клиентов при холодном старте
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	rekognitionClient := rekognition.NewFromConfig(cfg)
	imageAnalyzer = analyzer.NewImageAnalyzer(rekognitionClient)
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

	var analysisRequest types.ImageAnalysisRequestApi
	if err := json.Unmarshal([]byte(request.Body), &analysisRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"Invalid request body"}`,
			Headers:    headers,
		}, nil
	}

	// Проверка входных данных
	if analysisRequest.ImageURL == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"image_url is required"}`,
			Headers:    headers,
		}, nil
	}

	// Анализ изображения
	labels, err := imageAnalyzer.AnalyzeImage(ctx, analysisRequest)
	if err != nil {
		log.Printf("Failed to analyze image: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"Failed to analyze image"}`,
			Headers:    headers,
		}, nil
	}

	// Преобразование меток в категории
	var categories []string
	for _, label := range labels {
		if cats, ok := analyzer.LabelCategories[label]; ok {
			categories = append(categories, cats...)
		}
	}

	response := types.ApiResponse{
		Success: true,
		Data: types.ImageAnalysisResponseApi{
			Labels:     labels,
			Categories: categories,
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
