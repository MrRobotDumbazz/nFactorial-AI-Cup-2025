# Gift Recommendation Service

Сервис рекомендации подарков с использованием AWS Lambda и различных сервисов AWS.

## Структура API

### 1. Анализ изображений
- **Endpoint**: `POST /analyze-image`
- **Request**:
```json
{
    "image_url": "https://example.com/image.jpg"
}
```
- **Response**:
```json
{
    "success": true,
    "data": {
        "labels": ["phone", "electronics", "smartphone"],
        "categories": ["electronics", "gadgets"]
    }
}
```

### 2. Перевод текста
- **Endpoint**: `POST /translate`
- **Request**:
```json
{
    "text": "Hello world",
    "target_language": "ru"
}
```
- **Response**:
```json
{
    "success": true,
    "data": {
        "translated_text": "Привет мир"
    }
}
```

### 3. Озвучка текста
- **Endpoint**: `POST /text-to-speech`
- **Request**:
```json
{
    "text": "Hello world",
    "language": "en-US"
}
```
- **Response**:
```json
{
    "success": true,
    "data": {
        "audio_url": "https://your-bucket.s3.amazonaws.com/audio/123.mp3"
    }
}
```

### 4. Поиск товаров
- **Endpoint**: `POST /search-products`
- **Request**:
```json
{
    "categories": ["electronics", "gadgets"],
    "price_range": {
        "min": 100,
        "max": 1000
    },
    "marketplace": "amazon"
}
```
- **Response**:
```json
{
    "success": true,
    "data": {
        "products": [
            {
                "id": "123",
                "name": "Smartphone",
                "description": "Latest model",
                "price": 599.99,
                "category": "electronics",
                "image_url": "https://example.com/phone.jpg"
            }
        ]
    }
}
```

## Тестирование в Postman

1. Создайте новую коллекцию в Postman
2. Добавьте переменную окружения `api_url` с базовым URL вашего API Gateway
3. Импортируйте следующие запросы:

### Анализ изображений
```
POST {{api_url}}/analyze-image
Content-Type: application/json

{
    "image_url": "https://example.com/image.jpg"
}
```

### Перевод текста
```
POST {{api_url}}/translate
Content-Type: application/json

{
    "text": "Hello world",
    "target_language": "ru"
}
```

### Озвучка текста
```
POST {{api_url}}/text-to-speech
Content-Type: application/json

{
    "text": "Hello world",
    "language": "en-US"
}
```

### Поиск товаров
```
POST {{api_url}}/search-products
Content-Type: application/json

{
    "categories": ["electronics", "gadgets"],
    "price_range": {
        "min": 100,
        "max": 1000
    },
    "marketplace": "amazon"
}
```

## Переменные окружения

Для работы Lambda функций требуются следующие переменные окружения:

- `AUDIO_BUCKET_NAME` - имя S3 бакета для хранения аудио файлов
- `AWS_REGION` - регион AWS (например, us-east-1)
- `DYNAMODB_TABLE` - имя таблицы DynamoDB для хранения товаров

## Развертывание

1. Убедитесь, что у вас установлен AWS CLI и настроены учетные данные
2. Выполните `go mod tidy` для установки зависимостей
3. Соберите каждую Lambda функцию:
```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/image-analyzer/main.go
zip image-analyzer.zip bootstrap
# Повторите для остальных функций
```
4. Загрузите ZIP-файлы в AWS Lambda
5. Настройте API Gateway для маршрутизации запросов к соответствующим функциям 