package main

import (
	"context"
	"log"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/analyzer"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
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

func handleRequest(ctx context.Context, request types.ImageAnalysisRequest) (types.APIResponse, error) {
	// Проверка входных данных
	if request.ImageURL == "" {
		return types.APIResponse{
			Success: false,
			Error:   "image_url is required",
		}, nil
	}

	// Анализ изображения
	labels, err := imageAnalyzer.AnalyzeImage(ctx, request.ImageURL)
	if err != nil {
		log.Printf("Failed to analyze image: %v", err)
		return types.APIResponse{
			Success: false,
			Error:   "failed to analyze image",
		}, nil
	}

	// Преобразование меток в категории
	var categories []string
	for _, label := range labels {
		if cats, ok := analyzer.LabelCategories[label]; ok {
			categories = append(categories, cats...)
		}
	}

	return types.APIResponse{
		Success: true,
		Data: types.ImageAnalysisResponse{
			Labels:     labels,
			Categories: categories,
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
