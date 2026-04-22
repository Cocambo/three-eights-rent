package dto

import "time"

type RecommendationItemResponse struct {
	ID           uint     `json:"id"`
	Brand        string   `json:"brand"`
	Model        string   `json:"model"`
	Year         int      `json:"year"`
	FuelType     string   `json:"fuel_type"`
	Transmission string   `json:"transmission"`
	BodyType     string   `json:"body_type"`
	SeatsCount   int      `json:"seats_count"`
	PricePerDay  int64    `json:"price_per_day"`
	Purpose      string   `json:"purpose"`
	MainImageURL *string  `json:"main_image_url"`
	Score        *float64 `json:"score,omitempty"`
}

type RecommendationsResponse struct {
	Items []RecommendationItemResponse `json:"items"`
}

type RebuildRecommendationsResponse struct {
	InteractionsCount    int       `json:"interactions_count"`
	SourceCarsCount      int       `json:"source_cars_count"`
	RecommendationsCount int       `json:"recommendations_count"`
	RebuiltAt            time.Time `json:"rebuilt_at"`
}
