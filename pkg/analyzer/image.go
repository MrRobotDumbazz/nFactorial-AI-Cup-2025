package analyzer

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

type ImageAnalyzer struct {
	client *rekognition.Client
}

func NewImageAnalyzer(client *rekognition.Client) *ImageAnalyzer {
	return &ImageAnalyzer{client: client}
}

func (a *ImageAnalyzer) AnalyzeImage(ctx context.Context, imageURL string) ([]string, error) {
	// Загружаем изображение
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %v", err)
	}

	// Анализируем изображение с помощью Rekognition
	input := &rekognition.DetectLabelsInput{
		Image: &types.Image{
			Bytes: imageBytes,
		},
		MaxLabels:     aws.Int32(10),
		MinConfidence: aws.Float32(70.0),
	}

	output, err := a.client.DetectLabels(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze image: %v", err)
	}

	// Извлекаем метки
	var labels []string
	for _, label := range output.Labels {
		labels = append(labels, *label.Name)
	}

	return labels, nil
}

// Маппинг меток Rekognition на категории подарков
var LabelCategories = map[string][]string{
	"Sports":      {"sports"},
	"Electronics": {"electronics"},
	"Book":        {"books"},
	"Game":        {"toys", "electronics"},
	"Pet":         {"home"},
	"Music":       {"electronics"},
	"Art":         {"home"},
	"Food":        {"home"},
	"Fashion":     {"beauty"},
	"Technology":  {"electronics"},
	"Fitness":     {"sports"},
	"Baby":        {"toys"},
	"Garden":      {"home"},
	"Beauty":      {"beauty"},
}
