package main

import (
	"context"
	"log"
	"os"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/analyzer"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/translator"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/translate"
)

type Handler struct {
	imageAnalyzer *analyzer.ImageAnalyzer
	translator    *translator.Translator
	dynamoClient  *dynamodb.Client
}

func main() {
	// Загружаем конфигурацию AWS
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// Инициализируем клиенты AWS
	rekognitionClient := rekognition.NewFromConfig(cfg)
	translateClient := translate.NewFromConfig(cfg)
	pollyClient := polly.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)
	dynamoClient := dynamodb.NewFromConfig(cfg)

	// Создаем обработчик
	handler := &Handler{
		imageAnalyzer: analyzer.NewImageAnalyzer(rekognitionClient),
		translator: translator.NewTranslator(
			translateClient,
			pollyClient,
			s3Client,
			os.Getenv("AUDIO_BUCKET_NAME"),
		),
		dynamoClient: dynamoClient,
	}

	// Запускаем Lambda
	lambda.Start(handler.HandleRequest)
}

func (h *Handler) HandleRequest(ctx context.Context, request types.GiftRequest) (types.GiftRecommendation, error) {
	var categories []string

	// Анализируем изображение, если оно предоставлено
	if request.ImageURL != "" {
		labels, err := h.imageAnalyzer.AnalyzeImage(ctx, request.ImageURL)
		if err != nil {
			log.Printf("failed to analyze image: %v", err)
		} else {
			// Добавляем категории на основе меток изображения
			for _, label := range labels {
				if cats, ok := analyzer.LabelCategories[label]; ok {
					categories = append(categories, cats...)
				}
			}
		}
	}

	// Добавляем категории на основе возраста
	ageGroup := getAgeGroup(request.Age)
	if ageCats, ok := types.AgeCategories[ageGroup]; ok {
		categories = append(categories, ageCats...)
	}

	// Добавляем категории на основе повода
	if occasionCats, ok := types.OccasionCategories[request.Occasion]; ok {
		categories = append(categories, occasionCats...)
	}

	// Формируем рекомендации
	recommendations := h.getRecommendations(ctx, categories, request)

	// Переводим и озвучиваем описание, если требуется
	summary := generateSummary(recommendations, request)
	if request.Language != "en" {
		translatedSummary, err := h.translator.TranslateText(ctx, summary, request.Language)
		if err != nil {
			log.Printf("failed to translate summary: %v", err)
		} else {
			summary = translatedSummary
		}
	}

	result := types.GiftRecommendation{
		Products: recommendations,
		Summary:  summary,
	}

	// Создаем аудио-версию, если требуется
	if request.VoiceEnabled {
		audioURL, err := h.translator.TextToSpeech(ctx, summary, request.Language)
		if err != nil {
			log.Printf("failed to create audio: %v", err)
		} else {
			result.AudioURL = audioURL
		}
	}

	return result, nil
}

func (h *Handler) getRecommendations(ctx context.Context, categories []string, request types.GiftRequest) []types.Product {
	// TODO: Implement marketplace-specific product fetching
	// For now, return mock data
	return []types.Product{
		{
			Title:       "Sample Product 1",
			Description: "This is a sample product description",
			Price:       99.99,
			Rating:      4.5,
			URL:         "https://example.com/product1",
			ImageURL:    "https://example.com/product1.jpg",
			Store:       request.Marketplace,
			Category:    categories[0],
		},
	}
}

func generateSummary(products []types.Product, request types.GiftRequest) string {
	// TODO: Implement better summary generation
	return "Here are some gift recommendations based on your preferences..."
}

func getAgeGroup(age int) string {
	switch {
	case age <= 12:
		return "child"
	case age <= 19:
		return "teen"
	case age <= 59:
		return "adult"
	default:
		return "senior"
	}
}
