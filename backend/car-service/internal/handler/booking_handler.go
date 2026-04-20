package handler

import (
	"net/http"

	"car-service/internal/dto"
	"car-service/internal/middleware"
	"car-service/internal/service"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

func (h *BookingHandler) Create(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		writeError(c, serviceUnauthorizedError())
		return
	}

	var req dto.CreateBookingRequest
	if !bindJSON(c, &req) {
		return
	}

	booking, err := h.bookingService.CreateBooking(c.Request.Context(), service.CreateBookingCommand{
		UserID:    userID,
		CarID:     req.CarID,
		StartDate: *req.StartDate,
		EndDate:   *req.EndDate,
	})
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, http.StatusCreated, toBookingResponse(booking))
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		writeError(c, serviceUnauthorizedError())
		return
	}

	var req dto.BookingIDURI
	if !bindURI(c, &req) {
		return
	}

	booking, err := h.bookingService.CancelBooking(c.Request.Context(), userID, req.ID)
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, dto.CancelBookingResponse{
		BookingID: booking.ID,
		Status:    booking.Status,
		Message:   "booking cancelled",
	})
}

func (h *BookingHandler) List(c *gin.Context) {
	var req dto.ListBookingsQuery
	if !bindQuery(c, &req) {
		return
	}

	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		writeError(c, serviceUnauthorizedError())
		return
	}

	bookings, err := h.bookingService.GetUserBookings(c.Request.Context(), userID)
	if err != nil {
		writeError(c, err)
		return
	}

	writeSuccess(c, http.StatusOK, toListBookingsResponse(bookings))
}
