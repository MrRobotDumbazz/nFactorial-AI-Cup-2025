package main

import (
	"context"
	"log"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/translator"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/translate"
)

var translatorService *translator.Translator

func init() {
	// Инициализация AWS клиентов при холодном старте
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	translateClient := translate.NewFromConfig(cfg)
	translatorService = translator.NewTranslator(translateClient, nil, nil, "")
}

func handleRequest(ctx context.Context, request types.TranslationRequest) (types.APIResponse, error) {
	// Проверка входных данных
	if request.Text == "" {
		return types.APIResponse{
			Success: false,
			Error:   "text is required",
		}, nil
	}

	if request.TargetLanguage == "" {
		return types.APIResponse{
			Success: false,
			Error:   "target_language is required",
		}, nil
	}

	// Перевод текста
	translatedText, err := translatorService.TranslateText(ctx, request.Text, request.TargetLanguage)
	if err != nil {
		log.Printf("Failed to translate text: %v", err)
		return types.APIResponse{
			Success: false,
			Error:   "failed to translate text",
		}, nil
	}

	return types.APIResponse{
		Success: true,
		Data: types.TranslationResponse{
			TranslatedText: translatedText,
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
