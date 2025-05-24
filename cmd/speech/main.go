package main

import (
	"context"
	"log"
	"os"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/translator"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
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

func handleRequest(ctx context.Context, request types.SpeechRequest) (types.APIResponse, error) {
	// Проверка входных данных
	if request.Text == "" {
		return types.APIResponse{
			Success: false,
			Error:   "text is required",
		}, nil
	}

	if request.Language == "" {
		return types.APIResponse{
			Success: false,
			Error:   "language is required",
		}, nil
	}

	// Преобразование текста в речь
	audioURL, err := speechService.TextToSpeech(ctx, request.Text, request.Language)
	if err != nil {
		log.Printf("Failed to synthesize speech: %v", err)
		return types.APIResponse{
			Success: false,
			Error:   "failed to synthesize speech",
		}, nil
	}

	return types.APIResponse{
		Success: true,
		Data: types.SpeechResponse{
			AudioURL: audioURL,
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
