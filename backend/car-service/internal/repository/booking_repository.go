package repository

import (
	"context"
	"errors"
	"time"

	"car-service/internal/domains"
	apperrors "car-service/internal/errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type CreateBookingParams struct {
	UserID    uint
	CarID     uint
	StartDate time.Time
	EndDate   time.Time
}

type BookingInterval struct {
	StartDate time.Time `gorm:"column:start_date"`
	EndDate   time.Time `gorm:"column:end_date"`
}

type BookingRepository interface {
	Create(ctx context.Context, params CreateBookingParams) (domains.Booking, error)
	CancelByIDAndUser(ctx context.Context, bookingID, userID uint) (domains.Booking, error)
	ListByUserID(ctx context.Context, userID uint) ([]domains.Booking, error)
	ListActiveIntervalsByCarID(ctx context.Context, carID uint, fromDate, toDate time.Time) ([]BookingInterval, error)
}

type gormBookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &gormBookingRepository{db: db}
}

func (r *gormBookingRepository) Create(
	ctx context.Context,
	params CreateBookingParams,
) (domains.Booking, error) {
	booking := domains.Booking{
		UserID:    params.UserID,
		CarID:     params.CarID,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		Status:    domains.BookingStatusActive,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", int64(params.CarID)).Error; err != nil {
			return err
		}

		hasOverlap, err := hasActiveOverlap(tx, params.CarID, params.StartDate, params.EndDate)
		if err != nil {
			return err
		}

		if hasOverlap {
			return apperrors.New(apperrors.ErrConflict, "booking dates overlap with an existing booking")
		}

		return tx.Create(&booking).Error
	})
	if err != nil {
		return domains.Booking{}, mapBookingRepositoryError(err)
	}

	return booking, nil
}

func (r *gormBookingRepository) CancelByIDAndUser(
	ctx context.Context,
	bookingID, userID uint,
) (domains.Booking, error) {
	var booking domains.Booking

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("id = ? AND user_id = ?", bookingID, userID).
			Take(&booking).Error; err != nil {
			return err
		}

		switch booking.Status {
		case domains.BookingStatusCancelled:
			return nil
		case domains.BookingStatusCompleted:
			return apperrors.New(apperrors.ErrConflict, "completed booking cannot be cancelled")
		}

		cancelledAt := time.Now().UTC()
		if err := tx.Model(&domains.Booking{}).
			Where("id = ?", booking.ID).
			Updates(map[string]any{
				"status":       domains.BookingStatusCancelled,
				"cancelled_at": cancelledAt,
				"updated_at":   cancelledAt,
			}).Error; err != nil {
			return err
		}

		booking.Status = domains.BookingStatusCancelled
		booking.CancelledAt = &cancelledAt
		booking.UpdatedAt = cancelledAt

		return nil
	})
	if err != nil {
		return domains.Booking{}, mapBookingRepositoryError(err)
	}

	return booking, nil
}

func (r *gormBookingRepository) ListByUserID(ctx context.Context, userID uint) ([]domains.Booking, error) {
	var bookings []domains.Booking

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Preload("Car").
		Preload("Car.CarImages", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("is_main = ?", true).
				Order("is_main DESC").
				Order("sort_order ASC").
				Order("id ASC")
		}).
		Find(&bookings).Error
	if err != nil {
		return nil, mapBookingRepositoryError(err)
	}

	return bookings, nil
}

func (r *gormBookingRepository) ListActiveIntervalsByCarID(
	ctx context.Context,
	carID uint,
	fromDate, toDate time.Time,
) ([]BookingInterval, error) {
	var intervals []BookingInterval

	err := r.db.WithContext(ctx).
		Model(&domains.Booking{}).
		Select("start_date", "end_date").
		Where(
			"car_id = ? AND status = ? AND start_date < ? AND end_date > ?",
			carID,
			domains.BookingStatusActive,
			toDate,
			fromDate,
		).
		Order("start_date ASC").
		Order("end_date ASC").
		Find(&intervals).Error
	if err != nil {
		return nil, mapBookingRepositoryError(err)
	}

	return intervals, nil
}

func hasActiveOverlap(tx *gorm.DB, carID uint, startDate, endDate time.Time) (bool, error) {
	var booking domains.Booking

	err := tx.Model(&domains.Booking{}).
		Select("id").
		Where(
			"car_id = ? AND status = ? AND start_date < ? AND end_date > ?",
			carID,
			domains.BookingStatusActive,
			endDate,
			startDate,
		).
		Take(&booking).Error
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func mapBookingRepositoryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, apperrors.ErrConflict) {
		return err
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return apperrors.New(apperrors.ErrNotFound, "car not found")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return apperrors.New(apperrors.ErrNotFound, "car not found")
	}

	return mapRepositoryError(err, "booking not found")
}
