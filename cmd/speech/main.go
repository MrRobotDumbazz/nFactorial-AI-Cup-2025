package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/translator"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var speechService *translator.Translator

func init() {
	// Инициализация AWS клиентов при холодном старте
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	pollyClient := polly.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)

	// Получаем имя S3 бакета из переменных окружения
	bucketName := os.Getenv("AUDIO_BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("AUDIO_BUCKET_NAME environment variable is required")
	}

	speechService = translator.NewTranslator(nil, pollyClient, s3Client, bucketName)
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

	var speechRequest types.SpeechRequestApi
	if err := json.Unmarshal([]byte(request.Body), &speechRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"Invalid request body"}`,
			Headers:    headers,
		}, nil
	}

	if speechRequest.Text == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"text is required"}`,
			Headers:    headers,
		}, nil
	}

	if speechRequest.Language == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"language is required"}`,
			Headers:    headers,
		}, nil
	}

	// Преобразование текста в речь
	audioURL, err := speechService.TextToSpeech(ctx, speechRequest.Text, speechRequest.Language)
	if err != nil {
		log.Printf("Failed to synthesize speech: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"failed to synthesize speech"}`,
			Headers:    headers,
		}, nil
	}

	response := types.ApiResponse{
		Success: true,
		Data: types.SpeechResponseApi{
			AudioURL: audioURL,
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
