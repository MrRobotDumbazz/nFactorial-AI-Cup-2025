package main

import (
	"context"
	"encoding/json"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/analyzer"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/config"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var imageAnalyzer *analyzer.ImageAnalyzer

func init() {
	// Инициализируем AWS сервисы при старте Lambda
	ctx := context.Background()
	rekognitionClient, _, _, err := config.InitAWSServices(ctx)
	if err != nil {
		panic(err)
	}

	imageAnalyzer = analyzer.NewImageAnalyzer(rekognitionClient)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Для OPTIONS запросов сразу возвращаем CORS заголовки
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    types.GetCORSHeaders(),
		}, nil
	}

	// Парсим входящий JSON
	var analysisRequest types.ImageAnalysisRequestApi
	if err := json.Unmarshal([]byte(request.Body), &analysisRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"Invalid request body"}`,
			Headers:    types.GetCORSHeaders(),
		}, nil
	}

	// Анализируем изображение
	labels, err := imageAnalyzer.AnalyzeImage(ctx, analysisRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"Failed to analyze image"}`,
			Headers:    types.GetCORSHeaders(),
		}, nil
	}

	// Формируем ответ
	response := types.ImageAnalysisResponse{
		Labels: labels,
	}

	// Сериализуем ответ в JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"Failed to marshal response"}`,
			Headers:    types.GetCORSHeaders(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
		Headers:    types.GetCORSHeaders(),
	}, nil
}

func main() {
	// Запускаем Lambda
	lambda.Start(handleRequest)
}
