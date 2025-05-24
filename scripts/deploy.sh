#!/bin/bash

# Проверяем наличие переменных окружения
if [ -z "$AUDIO_BUCKET_NAME" ]; then
    echo "Error: AUDIO_BUCKET_NAME is not set"
    exit 1
fi

if [ -z "$AWS_ACCOUNT_ID" ]; then
    echo "Error: AWS_ACCOUNT_ID is not set"
    exit 1
fi

# Создаем S3 бакет, если он не существует
aws s3 mb s3://$AUDIO_BUCKET_NAME || true

# Создаем таблицу DynamoDB, если она не существует
aws dynamodb create-table \
    --table-name gift_recommendations \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 || true

# Создаем IAM роль
ROLE_NAME="gift-recommender-lambda-role"

# Заменяем переменную в IAM policy
sed "s/\${AUDIO_BUCKET_NAME}/$AUDIO_BUCKET_NAME/g" iam/lambda-role.json > /tmp/lambda-role.json

# Создаем роль и прикрепляем policy
aws iam create-role \
    --role-name $ROLE_NAME \
    --assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}' || true

aws iam put-role-policy \
    --role-name $ROLE_NAME \
    --policy-name "gift-recommender-policy" \
    --policy-document file:///tmp/lambda-role.json || true

# Ждем, пока роль будет создана
sleep 10

# Собираем и упаковываем Lambda функцию
GOOS=linux GOARCH=amd64 go build -o main ./cmd/lambda
zip function.zip main

# Создаем или обновляем Lambda функцию
FUNCTION_NAME="gift-recommender"
ROLE_ARN="arn:aws:iam::$AWS_ACCOUNT_ID:role/$ROLE_NAME"

if aws lambda get-function --function-name $FUNCTION_NAME 2>/dev/null; then
    # Обновляем существующую функцию
    aws lambda update-function-code \
        --function-name $FUNCTION_NAME \
        --zip-file fileb://function.zip
else
    # Создаем новую функцию
    aws lambda create-function \
        --function-name $FUNCTION_NAME \
        --runtime go1.x \
        --handler main \
        --role $ROLE_ARN \
        --zip-file fileb://function.zip
fi

# Очищаем временные файлы
rm main function.zip /tmp/lambda-role.json

echo "Deployment completed successfully!" 