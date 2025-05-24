package analyzer

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	customtypes "github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types" // Замените your-module на актуальное имя модуля

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

func (a *ImageAnalyzer) AnalyzeImage(ctx context.Context, request customtypes.ImageAnalysisRequest) ([]string, error) {
	var imageBytes []byte
	var err error

	switch request.ImageSource {
	case "url":
		if request.ImageURL == "" {
			return nil, fmt.Errorf("image_url is required for url source")
		}
		imageBytes, err = a.downloadImage(request.ImageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to download image: %v", err)
		}

	case "base64":
		if request.ImageBase64 == "" {
			return nil, fmt.Errorf("image_base64 is required for base64 source")
		}
		// Удаляем префикс data:image/...;base64, если он есть
		base64Data := request.ImageBase64
		if idx := strings.Index(base64Data, ","); idx != -1 {
			base64Data = base64Data[idx+1:]
		}
		imageBytes, err = base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 image: %v", err)
		}

	case "file":
		if len(request.ImageFile) == 0 {
			return nil, fmt.Errorf("image_file is required for file source")
		}
		imageBytes = request.ImageFile

	default:
		return nil, fmt.Errorf("invalid image source: %s", request.ImageSource)
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

func (a *ImageAnalyzer) downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
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
