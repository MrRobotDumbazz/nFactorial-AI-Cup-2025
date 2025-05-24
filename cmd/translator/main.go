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

	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
		}, nil
	}

	var translationRequest types.TranslationRequestApi
	if err := json.Unmarshal([]byte(request.Body), &translationRequest); err != nil {
		log.Printf("Failed to unmarshal request body: %v. Body: %s", err, request.Body)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error":"Invalid request body: %v"}`, err),
			Headers:    headers,
		}, nil
	}

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
