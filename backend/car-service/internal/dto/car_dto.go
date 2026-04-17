package dto

import "fmt"

const (
	DefaultCarsLimit = 20
	MaxCarsLimit     = 100
)

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
	SeatsMin     *int   `form:"seats_min" binding:"omitempty,gt=0"`
	PriceMin     *int64 `form:"price_min" binding:"omitempty,gte=0"`
	PriceMax     *int64 `form:"price_max" binding:"omitempty,gte=0"`
	Purpose      string `form:"purpose" binding:"omitempty,max=100"`
	SortBy       string `form:"sort_by" binding:"omitempty,oneof=id year price_per_day created_at"`
	SortOrder    string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	Limit        *int   `form:"limit" binding:"omitempty,gte=1,lte=100"`
	Offset       *int   `form:"offset" binding:"omitempty,gte=0"`
}

func (q *ListCarsQuery) Normalize() {
	if q.Limit == nil {
		limit := DefaultCarsLimit
		q.Limit = &limit
	}

	if q.Offset == nil {
		offset := 0
		q.Offset = &offset
	}

	if q.SortBy == "" {
		q.SortBy = "id"
	}

	if q.SortOrder == "" {
		q.SortOrder = "desc"
	}
}

func (q ListCarsQuery) Validate() error {
	if q.YearFrom != nil && q.YearTo != nil && *q.YearFrom > *q.YearTo {
		return fmt.Errorf("year_to must be greater than or equal to year_from")
	}

	if q.PriceMin != nil && q.PriceMax != nil && *q.PriceMin > *q.PriceMax {
		return fmt.Errorf("price_max must be greater than or equal to price_min")
	}

	return nil
}

func (q ListCarsQuery) LimitValue() int {
	if q.Limit == nil {
		return DefaultCarsLimit
	}

	return *q.Limit
}

func (q ListCarsQuery) OffsetValue() int {
	if q.Offset == nil {
		return 0
	}

	return *q.Offset
}

type GetCarByIDURI struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}

type CarImageResponse struct {
	ID        uint   `json:"id"`
	URL       string `json:"url"`
	IsMain    bool   `json:"is_main"`
	SortOrder int    `json:"sort_order"`
}

type CarCatalogItemResponse struct {
	ID           uint   `json:"id"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	FuelType     string `json:"fuel_type"`
	Transmission string `json:"transmission"`
	BodyType     string `json:"body_type"`
	SeatsCount   int    `json:"seats_count"`
	PricePerDay  int64  `json:"price_per_day"`
	Purpose      string `json:"purpose"`
	MainImageURL string `json:"main_image_url,omitempty"`
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
