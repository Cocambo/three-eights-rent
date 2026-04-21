package service

import (
	"context"
	"time"

	"car-service/internal/domains"
	apperrors "car-service/internal/errors"
	"car-service/internal/repository"
	"car-service/internal/storage"
)

type CreateBookingCommand struct {
	UserID    uint
	CarID     uint
	StartDate time.Time
	EndDate   time.Time
}

type BookingRecord struct {
	ID          uint
	CarID       uint
	StartDate   time.Time
	EndDate     time.Time
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CancelledAt *time.Time
}

type BookingInterval struct {
	StartDate time.Time
	EndDate   time.Time
}

type BookingHistoryItem struct {
	BookingRecord
	Car CatalogCar
}

type CarAvailabilityResult struct {
	CarID         uint
	From          time.Time
	To            time.Time
	BusyIntervals []BookingInterval
}

type BookingService interface {
	CreateBooking(ctx context.Context, command CreateBookingCommand) (BookingRecord, error)
	CancelBooking(ctx context.Context, userID, bookingID uint) (BookingRecord, error)
	GetUserBookings(ctx context.Context, userID uint) ([]BookingHistoryItem, error)
	GetCarAvailability(ctx context.Context, carID uint, fromDate, toDate string) (CarAvailabilityResult, error)
}

type bookingService struct {
	bookingRepository repository.BookingRepository
	carRepository     repository.CarRepository
	imageStorage      storage.ImageStorageService
}

func NewBookingService(
	bookingRepository repository.BookingRepository,
	carRepository repository.CarRepository,
	imageStorage storage.ImageStorageService,
) BookingService {
	return &bookingService{
		bookingRepository: bookingRepository,
		carRepository:     carRepository,
		imageStorage:      imageStorage,
	}
}

func (s *bookingService) CreateBooking(
	ctx context.Context,
	command CreateBookingCommand,
) (BookingRecord, error) {
	if command.UserID == 0 {
		return BookingRecord{}, validationError("user_id must be greater than zero")
	}

	if command.CarID == 0 {
		return BookingRecord{}, validationError("car_id must be greater than zero")
	}

	startDate, endDate, err := normalizeBookingDates(command.StartDate, command.EndDate)
	if err != nil {
		return BookingRecord{}, err
	}

	exists, err := s.carRepository.ExistsByID(ctx, command.CarID)
	if err != nil {
		return BookingRecord{}, err
	}

	if !exists {
		return BookingRecord{}, apperrors.New(apperrors.ErrNotFound, "car not found")
	}

	booking, err := s.bookingRepository.Create(ctx, repository.CreateBookingParams{
		UserID:    command.UserID,
		CarID:     command.CarID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return BookingRecord{}, err
	}

	return toBookingRecord(booking), nil
}

func (s *bookingService) CancelBooking(
	ctx context.Context,
	userID, bookingID uint,
) (BookingRecord, error) {
	if userID == 0 {
		return BookingRecord{}, validationError("user_id must be greater than zero")
	}

	if bookingID == 0 {
		return BookingRecord{}, validationError("booking_id must be greater than zero")
	}

	booking, err := s.bookingRepository.CancelByIDAndUser(ctx, bookingID, userID)
	if err != nil {
		return BookingRecord{}, err
	}

	return toBookingRecord(booking), nil
}

func (s *bookingService) GetUserBookings(
	ctx context.Context,
	userID uint,
) ([]BookingHistoryItem, error) {
	if userID == 0 {
		return nil, validationError("user_id must be greater than zero")
	}

	bookings, err := s.bookingRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	items := make([]BookingHistoryItem, 0, len(bookings))
	for i := range bookings {
		bookings[i].Car.CarImages = catalogImages(bookings[i].Car.CarImages)

		car, err := toCatalogCar(ctx, s.imageStorage, bookings[i].Car)
		if err != nil {
			return nil, err
		}

		items = append(items, BookingHistoryItem{
			BookingRecord: toBookingRecord(bookings[i]),
			Car:           car,
		})
	}

	return items, nil
}

func (s *bookingService) GetCarAvailability(
	ctx context.Context,
	carID uint,
	fromDate, toDate string,
) (CarAvailabilityResult, error) {
	if carID == 0 {
		return CarAvailabilityResult{}, validationError("car_id must be greater than zero")
	}

	startDate, err := parseDateOnly(fromDate, "from")
	if err != nil {
		return CarAvailabilityResult{}, err
	}

	endDate, err := parseDateOnly(toDate, "to")
	if err != nil {
		return CarAvailabilityResult{}, err
	}

	if !startDate.Before(endDate) {
		return CarAvailabilityResult{}, validationError("from must be earlier than to")
	}

	exists, err := s.carRepository.ExistsByID(ctx, carID)
	if err != nil {
		return CarAvailabilityResult{}, err
	}

	if !exists {
		return CarAvailabilityResult{}, apperrors.New(apperrors.ErrNotFound, "car not found")
	}

	intervals, err := s.bookingRepository.ListActiveIntervalsByCarID(ctx, carID, startDate, endDate)
	if err != nil {
		return CarAvailabilityResult{}, err
	}

	result := CarAvailabilityResult{
		CarID:         carID,
		From:          startDate,
		To:            endDate,
		BusyIntervals: make([]BookingInterval, 0, len(intervals)),
	}

	for _, interval := range intervals {
		result.BusyIntervals = append(result.BusyIntervals, BookingInterval{
			StartDate: interval.StartDate,
			EndDate:   interval.EndDate,
		})
	}

	return result, nil
}

func normalizeBookingDates(startDate, endDate time.Time) (time.Time, time.Time, error) {
	if startDate.IsZero() {
		return time.Time{}, time.Time{}, validationError("start_date is required")
	}

	if endDate.IsZero() {
		return time.Time{}, time.Time{}, validationError("end_date is required")
	}

	startDate = startDate.UTC()
	endDate = endDate.UTC()
	if !startDate.Before(endDate) {
		return time.Time{}, time.Time{}, validationError("start_date must be earlier than end_date")
	}

	return startDate, endDate, nil
}

func parseDateOnly(value, field string) (time.Time, error) {
	if value == "" {
		return time.Time{}, validationError(field + " is required")
	}

	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, validationError(field + " must be in YYYY-MM-DD format")
	}

	return parsed.UTC(), nil
}

func toBookingRecord(booking domains.Booking) BookingRecord {
	return BookingRecord{
		ID:          booking.ID,
		CarID:       booking.CarID,
		StartDate:   booking.StartDate,
		EndDate:     booking.EndDate,
		Status:      booking.Status,
		CreatedAt:   booking.CreatedAt,
		UpdatedAt:   booking.UpdatedAt,
		CancelledAt: booking.CancelledAt,
	}
}
