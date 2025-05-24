package analyzer

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
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

func (a *ImageAnalyzer) AnalyzeImage(ctx context.Context, request customtypes.ImageAnalysisRequestApi) ([]string, error) {
	log.Printf("Starting image analysis with source type: %s", request.ImageSource)

	var imageBytes []byte
	var err error

	switch request.ImageSource {
	case "url":
		if request.ImageURL == "" {
			return nil, fmt.Errorf("image_url is required for url source")
		}
		log.Printf("Downloading image from URL: %s", request.ImageURL)
		imageBytes, err = a.downloadImage(request.ImageURL)
		if err != nil {
			log.Printf("Failed to download image: %v", err)
			return nil, fmt.Errorf("failed to download image: %v", err)
		}
		log.Printf("Successfully downloaded image, size: %d bytes", len(imageBytes))

	case "base64":
		if request.ImageBase64 == "" {
			return nil, fmt.Errorf("image_base64 is required for base64 source")
		}
		log.Printf("Decoding base64 image")
		// Удаляем префикс data:image/...;base64, если он есть
		base64Data := request.ImageBase64
		if idx := strings.Index(base64Data, ","); idx != -1 {
			base64Data = base64Data[idx+1:]
		}
		imageBytes, err = base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			log.Printf("Failed to decode base64 image: %v", err)
			return nil, fmt.Errorf("failed to decode base64 image: %v", err)
		}
		log.Printf("Successfully decoded base64 image, size: %d bytes", len(imageBytes))

	case "file":
		if len(request.ImageFile) == 0 {
			return nil, fmt.Errorf("image_file is required for file source")
		}
		imageBytes = request.ImageFile
		log.Printf("Using provided file bytes, size: %d bytes", len(imageBytes))

	default:
		return nil, fmt.Errorf("invalid image source: %s", request.ImageSource)
	}

	// Анализируем изображение с помощью Rekognition
	log.Printf("Calling AWS Rekognition DetectLabels")
	input := &rekognition.DetectLabelsInput{
		Image: &types.Image{
			Bytes: imageBytes,
		},
		MaxLabels:     aws.Int32(10),
		MinConfidence: aws.Float32(70.0),
	}

	output, err := a.client.DetectLabels(ctx, input)
	if err != nil {
		log.Printf("AWS Rekognition DetectLabels failed: %v", err)
		return nil, fmt.Errorf("failed to analyze image: %v", err)
	}

	// Извлекаем метки
	var labels []string
	for _, label := range output.Labels {
		labels = append(labels, *label.Name)
	}
	log.Printf("AWS Rekognition detected %d labels: %v", len(labels), labels)

	return labels, nil
}

func (a *ImageAnalyzer) downloadImage(url string) ([]byte, error) {
	log.Printf("Starting image download from URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP GET request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		log.Printf("Download failed: %v", err)
		return nil, err
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}

	log.Printf("Successfully downloaded image, size: %d bytes", len(imageBytes))
	return imageBytes, nil
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
