package dto

import "time"

type HealthResponse struct {
	Status string `json:"status"`
}

type PaginationMeta struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

type ListCarsQuery struct {
	Q            string `form:"q" binding:"omitempty,max=100"`
	Brand        string `form:"brand" binding:"omitempty,max=100"`
	Model        string `form:"model" binding:"omitempty,max=100"`
	YearFrom     *int   `form:"year_from" binding:"omitempty,gte=1886,lte=2100"`
	YearTo       *int   `form:"year_to" binding:"omitempty,gte=1886,lte=2100"`
	FuelType     string `form:"fuel_type" binding:"omitempty,max=50"`
	Transmission string `form:"transmission" binding:"omitempty,max=50"`
	BodyType     string `form:"body_type" binding:"omitempty,max=50"`
	SeatsCount   *int   `form:"seats_count" binding:"omitempty,gt=0"`
	PriceMin     *int64 `form:"price_min" binding:"omitempty,gte=0"`
	PriceMax     *int64 `form:"price_max" binding:"omitempty,gte=0"`
	Purpose      string `form:"purpose" binding:"omitempty,max=100"`
	SortBy       string `form:"sort_by" binding:"omitempty,oneof=id year price_per_day created_at"`
	SortOrder    string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	Limit        *int   `form:"limit" binding:"omitempty,gte=1"`
	Offset       *int   `form:"offset" binding:"omitempty,gte=0"`
}

type GetCarByIDURI struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}

type UploadCarImageForm struct {
	IsMain    bool `form:"is_main"`
	SortOrder int  `form:"sort_order" binding:"omitempty,gte=0"`
}

type CarImageResponse struct {
	ID        uint   `json:"id"`
	URL       string `json:"url"`
	IsMain    bool   `json:"is_main"`
	SortOrder int    `json:"sort_order"`
}

type UploadedCarImageResponse struct {
	ID          uint      `json:"id"`
	CarID       uint      `json:"car_id"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type"`
	FileSize    int64     `json:"file_size"`
	IsMain      bool      `json:"is_main"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

type CarCatalogItemResponse struct {
	ID           uint    `json:"id"`
	Brand        string  `json:"brand"`
	Model        string  `json:"model"`
	Year         int     `json:"year"`
	FuelType     string  `json:"fuel_type"`
	Transmission string  `json:"transmission"`
	BodyType     string  `json:"body_type"`
	SeatsCount   int     `json:"seats_count"`
	PricePerDay  int64   `json:"price_per_day"`
	Purpose      string  `json:"purpose"`
	MainImageURL *string `json:"main_image_url"`
}

type CarDetailsResponse struct {
	ID           uint               `json:"id"`
	Brand        string             `json:"brand"`
	Model        string             `json:"model"`
	Year         int                `json:"year"`
	FuelType     string             `json:"fuel_type"`
	Transmission string             `json:"transmission"`
	BodyType     string             `json:"body_type"`
	Color        string             `json:"color"`
	SeatsCount   int                `json:"seats_count"`
	PricePerDay  int64              `json:"price_per_day"`
	Purpose      string             `json:"purpose"`
	Description  string             `json:"description"`
	Images       []CarImageResponse `json:"images"`
}

type CarsCatalogResponse struct {
	Items      []CarCatalogItemResponse `json:"items"`
	Pagination PaginationMeta           `json:"pagination"`
}
