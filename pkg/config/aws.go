package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/translate"
)

// InitAWSServices инициализирует AWS сервисы с стандартной конфигурацией
func InitAWSServices(ctx context.Context) (*rekognition.Client, *translate.Client, *polly.Client, error) {
	// Загружаем конфигурацию AWS из переменных окружения или файла credentials
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	// Создаем клиентов для каждого сервиса
	rekognitionClient := rekognition.NewFromConfig(cfg)
	translateClient := translate.NewFromConfig(cfg)
	pollyClient := polly.NewFromConfig(cfg)

	return rekognitionClient, translateClient, pollyClient, nil
}
