package repository

import (
	"context"
	"time"

	"car-service/internal/domains"

	"gorm.io/gorm"
)

const (
	favoriteInteractionWeight = 1.0
	bookingInteractionWeight  = 3.0
)

type RecommendationInteraction struct {
	UserID uint    `gorm:"column:user_id"`
	CarID  uint    `gorm:"column:car_id"`
	Weight float64 `gorm:"column:weight"`
}

type PopularCarScore struct {
	CarID uint    `gorm:"column:car_id"`
	Score float64 `gorm:"column:score"`
}

type RecommendationRepository interface {
	ListAggregatedInteractions(ctx context.Context) ([]RecommendationInteraction, error)
	ListUserAggregatedInteractions(ctx context.Context, userID uint) ([]RecommendationInteraction, error)
	ListBySourceCarIDs(ctx context.Context, sourceCarIDs []uint) ([]domains.CarRecommendation, error)
	ReplaceAll(ctx context.Context, recommendations []domains.CarRecommendation) error
	ListPopularCarScores(ctx context.Context, limit int) ([]PopularCarScore, error)
}

type gormRecommendationRepository struct {
	db *gorm.DB
}

func NewRecommendationRepository(db *gorm.DB) RecommendationRepository {
	return &gormRecommendationRepository{db: db}
}

func (r *gormRecommendationRepository) ListAggregatedInteractions(
	ctx context.Context,
) ([]RecommendationInteraction, error) {
	return r.scanInteractions(ctx, "")
}

func (r *gormRecommendationRepository) ListUserAggregatedInteractions(
	ctx context.Context,
	userID uint,
) ([]RecommendationInteraction, error) {
	return r.scanInteractions(ctx, "WHERE interactions.user_id = ?", userID)
}

func (r *gormRecommendationRepository) ListBySourceCarIDs(
	ctx context.Context,
	sourceCarIDs []uint,
) ([]domains.CarRecommendation, error) {
	if len(sourceCarIDs) == 0 {
		return []domains.CarRecommendation{}, nil
	}

	var recommendations []domains.CarRecommendation

	err := r.db.WithContext(ctx).
		Where("source_car_id IN ?", sourceCarIDs).
		Order("source_car_id ASC").
		Order("rank ASC").
		Find(&recommendations).Error
	if err != nil {
		return nil, mapRepositoryError(err, "car recommendations not found")
	}

	return recommendations, nil
}

func (r *gormRecommendationRepository) ReplaceAll(
	ctx context.Context,
	recommendations []domains.CarRecommendation,
) error {
	now := time.Now().UTC()
	for i := range recommendations {
		recommendations[i].CreatedAt = now
		recommendations[i].UpdatedAt = now
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).
			Delete(&domains.CarRecommendation{}).Error; err != nil {
			return err
		}

		if len(recommendations) == 0 {
			return nil
		}

		return tx.CreateInBatches(recommendations, 500).Error
	})
	if err != nil {
		return mapRepositoryError(err, "car recommendations not saved")
	}

	return nil
}

func (r *gormRecommendationRepository) ListPopularCarScores(
	ctx context.Context,
	limit int,
) ([]PopularCarScore, error) {
	if limit <= 0 {
		return []PopularCarScore{}, nil
	}

	var rows []PopularCarScore

	err := r.db.WithContext(ctx).
		Model(&domains.Booking{}).
		Select("bookings.car_id AS car_id, COUNT(*)::double precision AS score").
		Where("bookings.status = ?", domains.BookingStatusActive).
		Group("bookings.car_id").
		Order("score DESC").
		Order("bookings.car_id ASC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, mapRepositoryError(err, "popular cars not found")
	}

	return rows, nil
}

func aggregatedInteractionsQuery(outerFilter string) string {
	return `
		WITH interactions AS (
			SELECT favorites.user_id, favorites.car_id, ?::double precision AS weight
			FROM favorites

			UNION ALL

			SELECT booking_pairs.user_id, booking_pairs.car_id, ?::double precision AS weight
			FROM (
				SELECT DISTINCT bookings.user_id, bookings.car_id
				FROM bookings
				WHERE bookings.status = ?
			) AS booking_pairs
		)
		SELECT interactions.user_id, interactions.car_id, SUM(interactions.weight) AS weight
		FROM interactions
		` + outerFilter + `
		GROUP BY interactions.user_id, interactions.car_id
		ORDER BY interactions.user_id ASC, interactions.car_id ASC
	`
}

func (r *gormRecommendationRepository) scanInteractions(
	ctx context.Context,
	filter string,
	args ...any,
) ([]RecommendationInteraction, error) {
	queryArgs := []any{
		favoriteInteractionWeight,
		bookingInteractionWeight,
		domains.BookingStatusActive,
	}
	queryArgs = append(queryArgs, args...)

	var interactions []RecommendationInteraction

	err := r.db.WithContext(ctx).
		Raw(aggregatedInteractionsQuery(filter), queryArgs...).
		Scan(&interactions).Error
	if err != nil {
		return nil, mapRepositoryError(err, "recommendation interactions not found")
	}

	return interactions, nil
}
