package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/translator"
	"github.com/MrRobotDumbazz/nFactorial-AI-Cup-2025/pkg/types"
	"github.com/aws/aws-lambda-go/events"
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

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
		"Content-Type":                "application/json",
	}

	// Подробное логирование запроса
	log.Printf("=== START REQUEST DEBUGGING ===")
	log.Printf("Full request struct: %#v", request)
	log.Printf("Request body (raw): %q", request.Body)
	log.Printf("Request body length: %d", len(request.Body))
	log.Printf("Request headers: %#v", request.Headers)
	log.Printf("Request method: %s", request.HTTPMethod)
	log.Printf("Request path: %s", request.Path)
	log.Printf("IsBase64Encoded: %v", request.IsBase64Encoded)
	log.Printf("=== END REQUEST DEBUGGING ===")

	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
		}, nil
	}

	// Проверяем, что тело запроса не пустое
	if request.Body == "" {
		log.Printf("Empty request body received")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"Request body is empty"}`,
			Headers:    headers,
		}, nil
	}

	// Пробуем распарсить тело запроса напрямую
	var translationRequest types.TranslationRequestApi
	if err := json.Unmarshal([]byte(request.Body), &translationRequest); err != nil {
		log.Printf("Failed to unmarshal request body: %v. Body: %s", err, request.Body)

		// Если не получилось, пробуем распарсить как строку JSON
		var jsonStr string
		if err := json.Unmarshal([]byte(request.Body), &jsonStr); err == nil {
			// Если это строка JSON, пробуем распарсить её
			if err := json.Unmarshal([]byte(jsonStr), &translationRequest); err != nil {
				log.Printf("Failed to unmarshal JSON string: %v. String: %s", err, jsonStr)
				return events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       fmt.Sprintf(`{"error":"Invalid request body format: %v"}`, err),
					Headers:    headers,
				}, nil
			}
		} else {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf(`{"error":"Invalid request body: %v"}`, err),
				Headers:    headers,
			}, nil
		}
	}

	// Логируем успешно распарсенный запрос
	log.Printf("Successfully parsed translation request: %#v", translationRequest)

	if translationRequest.Text == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"text is required"}`,
			Headers:    headers,
		}, nil
	}

	if translationRequest.TargetLanguage == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error":"target_language is required"}`,
			Headers:    headers,
		}, nil
	}

	translatedText, err := translatorService.TranslateText(ctx, translationRequest.Text, translationRequest.TargetLanguage)
	if err != nil {
		log.Printf("Failed to translate text: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"failed to translate text"}`,
			Headers:    headers,
		}, nil
	}

	response := types.ApiResponse{
		Success: true,
		Data: types.TranslationResponseApi{
			TranslatedText: translatedText,
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
