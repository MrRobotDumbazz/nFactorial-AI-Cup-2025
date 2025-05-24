package types

// Общие типы для всех API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Анализ изображений
type ImageAnalysisRequest struct {
	ImageURL string `json:"image_url"`
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
