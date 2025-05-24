package types

type GiftRequest struct {
	Occasion     string   `json:"occasion"`      // Повод для подарка
	Gender       string   `json:"gender"`        // Пол получателя
	Age          int      `json:"age"`           // Возраст получателя
	Interests    []string `json:"interests"`     // Интересы
	PriceRange   Range    `json:"price_range"`   // Ценовой диапазон
	Marketplace  string   `json:"marketplace"`   // Предпочитаемый маркетплейс
	Language     string   `json:"language"`      // Язык ответа (ru/en/kk)
	ImageURL     string   `json:"image_url"`     // URL фото для анализа (опционально)
	VoiceEnabled bool     `json:"voice_enabled"` // Нужен ли голосовой ответ
}

type Range struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type GiftRecommendation struct {
	Products []Product `json:"products"`
	Summary  string    `json:"summary"`   // Текстовое описание рекомендаций
	AudioURL string    `json:"audio_url"` // URL аудио-версии (если запрошено)
}

type Product struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Rating      float64 `json:"rating"`
	URL         string  `json:"url"`
	ImageURL    string  `json:"image_url"`
	Store       string  `json:"store"`
	Category    string  `json:"category"`
}

// Маппинг категорий для разных маркетплейсов
var CategoryMappings = map[string]map[string]string{
	"kaspi": {
		"electronics": "Электроника",
		"books":       "Книги",
		"sports":      "Спорт и отдых",
		"beauty":      "Красота и здоровье",
		"toys":        "Детские товары",
		"home":        "Товары для дома",
	},
	"aliexpress": {
		"electronics": "Electronics",
		"books":       "Books & Office",
		"sports":      "Sports & Entertainment",
		"beauty":      "Beauty & Health",
		"toys":        "Toys & Hobbies",
		"home":        "Home & Garden",
	},
	"wildberries": {
		"electronics": "Электроника",
		"books":       "Книги",
		"sports":      "Спорт",
		"beauty":      "Красота",
		"toys":        "Детям",
		"home":        "Дом",
	},
	"ozon": {
		"electronics": "Электроника",
		"books":       "Книги",
		"sports":      "Спорт и отдых",
		"beauty":      "Красота и здоровье",
		"toys":        "Детские товары",
		"home":        "Дом и сад",
	},
}

// Маппинг поводов для подарков на категории
var OccasionCategories = map[string][]string{
	"birthday": {
		"electronics",
		"beauty",
		"sports",
		"home",
	},
	"wedding": {
		"home",
		"electronics",
	},
	"graduation": {
		"electronics",
		"books",
		"sports",
	},
	"newborn": {
		"toys",
		"home",
	},
}

// Маппинг возрастных групп на категории
var AgeCategories = map[string][]string{
	"child": { // 0-12
		"toys",
		"books",
		"sports",
	},
	"teen": { // 13-19
		"electronics",
		"sports",
		"books",
	},
	"adult": { // 20-59
		"electronics",
		"beauty",
		"home",
		"sports",
	},
	"senior": { // 60+
		"home",
		"books",
		"health",
	},
}
