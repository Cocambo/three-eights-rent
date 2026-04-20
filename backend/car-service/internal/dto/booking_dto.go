package dto

import "time"

type ListBookingsQuery struct{}

type BookingIDURI struct {
	ID uint `uri:"id" binding:"required,gt=0"`
}

type CreateBookingRequest struct {
	CarID     uint       `json:"car_id" binding:"required,gt=0"`
	StartDate *time.Time `json:"start_date" binding:"required"`
	EndDate   *time.Time `json:"end_date" binding:"required"`
}

type BookingResponse struct {
	ID          uint       `json:"id"`
	CarID       uint       `json:"car_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
}

type BookingHistoryItemResponse struct {
	ID          uint                   `json:"id"`
	CarID       uint                   `json:"car_id"`
	StartDate   time.Time              `json:"start_date"`
	EndDate     time.Time              `json:"end_date"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CancelledAt *time.Time             `json:"cancelled_at,omitempty"`
	Car         CarCatalogItemResponse `json:"car"`
}

type ListBookingsResponse struct {
	Items []BookingHistoryItemResponse `json:"items"`
}

type CancelBookingResponse struct {
	BookingID uint   `json:"booking_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
