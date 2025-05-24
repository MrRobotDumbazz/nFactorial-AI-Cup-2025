# nFactorial-AI-Cup-2025

Fork this repository and build nFactorial AI Cup 2025 projects

**Name:** MrRobotDumbazz

**App Name:** AI Gift Advisor

## Description

AI Gift Advisor - это интеллектуальный сервис рекомендации подарков, использующий различные AWS сервисы для анализа и персонализации подарков.

### Технические особенности

1. **Архитектура:**
   - Serverless архитектура на AWS Lambda
   - API Gateway для REST API
   - DynamoDB для хранения данных о товарах
   - S3 для хранения аудио файлов

2. **AWS Сервисы:**
   - Amazon Rekognition для анализа изображений
   - Amazon Translate для перевода описаний
   - Amazon Polly для озвучки описаний
   - DynamoDB для поиска товаров

3. **API Endpoints:**
   - `POST /analyze-image` - анализ изображений для определения категорий
   - `POST /translate` - перевод описаний товаров
   - `POST /text-to-speech` - озвучка описаний
   - `POST /search-products` - поиск товаров по категориям

4. **Стек технологий:**
   - Go 1.24.2
   - AWS SDK v2
   - AWS Lambda Go Runtime
   - DynamoDB
   - API Gateway

## Локальная разработка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/your-username/nFactorial-AI-Cup-2025.git
cd nFactorial-AI-Cup-2025
```

2. Установите зависимости:
```bash
go mod tidy
```

3. Настройте переменные окружения:
```bash
cp .env.example .env
# Отредактируйте .env файл
```

4. Соберите Lambda функции:
```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/image-analyzer/main.go
zip image-analyzer.zip bootstrap
# Повторите для остальных функций
```

## Развертывание

1. Создайте необходимые AWS ресурсы:
   - DynamoDB таблицу
   - S3 бакет
   - IAM роли (см. `/iam/README.md`)

2. Загрузите Lambda функции:
   - Создайте новые Lambda функции в AWS Console
   - Загрузите ZIP-файлы
   - Настройте триггеры API Gateway

3. Настройте переменные окружения в AWS Lambda:
   - `AUDIO_BUCKET_NAME`
   - `AWS_REGION`
   - `DYNAMODB_TABLE`

## Тестирование

Подробные инструкции по тестированию API через Postman находятся в разделе "Тестирование в Postman" ниже.

[Полная документация по тестированию API]

## Typeform для сдачи проекта

https://docs.google.com/forms/d/e/1FAIpQLSdjbTZXt-8P0OTyMEDTQDszE-YGI5KcLYsN6pwxHmX0Fa3tzg/viewform?usp=dialog

**DEADLINE:** 25/05/2025 10:00 