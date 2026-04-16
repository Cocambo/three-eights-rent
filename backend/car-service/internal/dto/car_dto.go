package dto

type CarResponse struct {
	ID             uint   `json:"id"`
	Brand          string `json:"brand"`
	Model          string `json:"model"`
	Year           int    `json:"year"`
	RegistrationNo string `json:"registration_no"`
	ImageURL       string `json:"image_url,omitempty"`
}

type HealthResponse struct {
	Status string `json:"status"`
}
