package translator

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	pollyTypes "github.com/aws/aws-sdk-go-v2/service/polly/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/translate"
)

type Translator struct {
	translateClient *translate.Client
	pollyClient     *polly.Client
	s3Client        *s3.Client
	bucketName      string
}

func NewTranslator(translateClient *translate.Client, pollyClient *polly.Client, s3Client *s3.Client, bucketName string) *Translator {
	return &Translator{
		translateClient: translateClient,
		pollyClient:     pollyClient,
		s3Client:        s3Client,
		bucketName:      bucketName,
	}
}

func (t *Translator) TranslateText(ctx context.Context, text, targetLang string) (string, error) {
	input := &translate.TranslateTextInput{
		Text:               aws.String(text),
		SourceLanguageCode: aws.String("auto"),
		TargetLanguageCode: aws.String(targetLang),
	}

	output, err := t.translateClient.TranslateText(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to translate text: %v", err)
	}

	return *output.TranslatedText, nil
}

func (t *Translator) TextToSpeech(ctx context.Context, text, lang string) (string, error) {
	// Выбираем голос в зависимости от языка
	voice := t.selectVoice(lang)

	// Конвертируем текст в речь
	input := &polly.SynthesizeSpeechInput{
		OutputFormat: pollyTypes.OutputFormatMp3,
		Text:         aws.String(text),
		VoiceId:      voice,
		Engine:       pollyTypes.EngineNeural,
	}

	output, err := t.pollyClient.SynthesizeSpeech(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to synthesize speech: %v", err)
	}

	// Генерируем уникальное имя файла (берем первые 20 символов текста)
	textPrefix := text
	if len(text) > 20 {
		textPrefix = text[:20]
	}
	key := fmt.Sprintf("audio/%s_%s.mp3", voice, strings.ReplaceAll(textPrefix, " ", "_"))

	// Загружаем аудио в S3
	_, err = t.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(t.bucketName),
		Key:    aws.String(key),
		Body:   output.AudioStream,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload audio to S3: %v", err)
	}

	// Возвращаем URL аудио файла
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", t.bucketName, key), nil
}

func (t *Translator) selectVoice(lang string) pollyTypes.VoiceId {
	switch lang {
	case "ru":
		return pollyTypes.VoiceIdMaxim
	case "kk":
		return pollyTypes.VoiceIdSalli // Используем английский голос, так как казахского нет
	default:
		return pollyTypes.VoiceIdJoanna // Английский по умолчанию
	}
}
