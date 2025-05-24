package types

// Общие типы для всех API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Анализ изображений
type ImageAnalysisRequest struct {
	ImageURL    string `json:"image_url,omitempty"`    // URL изображения
	ImageBase64 string `json:"image_base64,omitempty"` // Base64 encoded изображение
	ImageFile   []byte `json:"image_file,omitempty"`   // Бинарные данные файла
	ImageSource string `json:"image_source"`           // Тип источника: "url", "base64", "file"
}

type ImageAnalysisResponse struct {
	Labels     []string `json:"labels"`
	Categories []string `json:"categories"`
}

// Перевод текста
type TranslationRequest struct {
	Text           string `json:"text"`
	TargetLanguage string `json:"target_language"`
}

type TranslationResponse struct {
	TranslatedText string `json:"translated_text"`
}

// Озвучка текста
type SpeechRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

type SpeechResponse struct {
	AudioURL string `json:"audio_url"`
}

// Поиск товаров
type ProductSearchRequest struct {
	Categories  []string `json:"categories"`
	PriceRange  Range    `json:"price_range"`
	Marketplace string   `json:"marketplace"`
}

type ProductSearchResponse struct {
	Products []Product `json:"products"`
}

// Структуры для товаров
type RangeApi struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type ProductApi struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Category    string   `json:"category"`
	ImageURL    string   `json:"image_url"`
	Tags        []string `json:"tags,omitempty"`
}

// GetCORSHeaders возвращает стандартные CORS заголовки
func GetCORSHeaders() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Methods":     "GET,POST,PUT,DELETE,OPTIONS",
		"Access-Control-Allow-Headers":     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Access-Control-Allow-Credentials": "true",
		"Content-Type":                     "application/json",
	}
}
