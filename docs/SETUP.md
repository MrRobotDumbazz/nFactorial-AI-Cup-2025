# Инструкция по настройке и тестированию

## Структура Lambda функций

1. **gift-recommender** - основная Lambda функция для рекомендаций подарков
   - Путь: `/recommend`
   - Метод: POST
   - Входные данные: GiftRequest
   - Выходные данные: GiftRecommendation

## Настройка API Gateway

1. Создайте новый REST API:
```bash
aws apigateway create-rest-api \
    --name "gift-recommender-api" \
    --description "API для сервиса рекомендации подарков"
```

2. Получите root-id созданного API:
```bash
aws apigateway get-resources \
    --rest-api-id YOUR_API_ID
```

3. Создайте ресурс /recommend:
```bash
aws apigateway create-resource \
    --rest-api-id YOUR_API_ID \
    --parent-id ROOT_ID \
    --path-part "recommend"
```

4. Настройте метод POST:
```bash
aws apigateway put-method \
    --rest-api-id YOUR_API_ID \
    --resource-id RESOURCE_ID \
    --http-method POST \
    --authorization-type NONE
```

5. Интегрируйте с Lambda:
```bash
aws apigateway put-integration \
    --rest-api-id YOUR_API_ID \
    --resource-id RESOURCE_ID \
    --http-method POST \
    --type AWS_PROXY \
    --integration-http-method POST \
    --uri arn:aws:apigateway:REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:REGION:ACCOUNT_ID:function:gift-recommender/invocations
```

6. Разверните API:
```bash
aws apigateway create-deployment \
    --rest-api-id YOUR_API_ID \
    --stage-name prod
```

## Тестирование через Postman

1. **Endpoint**: 
```
https://YOUR_API_ID.execute-api.REGION.amazonaws.com/prod/recommend
```

2. **Метод**: POST

3. **Headers**:
```
Content-Type: application/json
```

4. **Body** (пример запроса):
```json
{
    "occasion": "birthday",
    "gender": "female",
    "age": 25,
    "interests": ["technology", "books"],
    "price_range": {
        "min": 5000,
        "max": 50000
    },
    "marketplace": "kaspi",
    "language": "ru",
    "image_url": "https://example.com/interests.jpg",
    "voice_enabled": true
}
```

5. **Пример ответа**:
```json
{
    "products": [
        {
            "title": "Kindle Paperwhite",
            "description": "Электронная книга с подсветкой",
            "price": 89990,
            "rating": 4.8,
            "url": "https://kaspi.kz/shop/kindle-paperwhite",
            "image_url": "https://resources.kaspi.kz/kindle.jpg",
            "store": "kaspi",
            "category": "electronics"
        }
    ],
    "summary": "Учитывая ваши интересы к технологиям и книгам, мы подобрали несколько отличных вариантов подарка...",
    "audio_url": "https://your-bucket.s3.amazonaws.com/audio/maxim_summary.mp3"
}
```

## Коды ответов

- 200: Успешный запрос
- 400: Некорректные входные данные
- 500: Внутренняя ошибка сервера

## Примеры тестовых сценариев

1. **День рождения для подростка**:
```json
{
    "occasion": "birthday",
    "gender": "male",
    "age": 15,
    "interests": ["gaming", "sports"],
    "price_range": {
        "min": 10000,
        "max": 100000
    },
    "marketplace": "wildberries",
    "language": "ru",
    "voice_enabled": false
}
```

2. **Свадебный подарок**:
```json
{
    "occasion": "wedding",
    "gender": "couple",
    "age": 30,
    "interests": ["home", "cooking"],
    "price_range": {
        "min": 50000,
        "max": 200000
    },
    "marketplace": "kaspi",
    "language": "kk",
    "voice_enabled": true
}
```

3. **Подарок для новорожденного**:
```json
{
    "occasion": "newborn",
    "gender": "female",
    "age": 0,
    "interests": ["toys", "care"],
    "price_range": {
        "min": 5000,
        "max": 30000
    },
    "marketplace": "ozon",
    "language": "ru",
    "voice_enabled": false
}
```

## Мониторинг и логи

1. **CloudWatch Logs**:
```bash
aws logs get-log-events \
    --log-group-name /aws/lambda/gift-recommender \
    --log-stream-name STREAM_NAME
```

2. **Метрики API Gateway**:
- Latency
- Integration Latency
- Requests
- 4XXError
- 5XXError

## Полезные команды

1. **Обновление Lambda функции**:
```bash
./scripts/deploy.sh
```

2. **Очистка CloudWatch логов**:
```bash
aws logs delete-log-group \
    --log-group-name /aws/lambda/gift-recommender
```

3. **Получение URL API**:
```bash
aws apigateway get-stages \
    --rest-api-id YOUR_API_ID
``` 