package dto

type CarResponse struct {
	ID           uint   `json:"id"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	FuelType     string `json:"fuel_type"`
	Transmission string `json:"transmission"`
	BodyType     string `json:"body_type"`
	Color        string `json:"color"`
	SeatsCount   int    `json:"seats_count"`
	PricePerDay  int64  `json:"price_per_day"`
	Purpose      string `json:"purpose"`
	Description  string `json:"description"`
}

type HealthResponse struct {
	Status string `json:"status"`
}
